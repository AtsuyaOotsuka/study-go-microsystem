package main

import (
	"context"
	"fmt"
	"io"
	"microservices/chat/internal/model"
	"microservices/chat/tests/test_funcs"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testServer *http.Server
var baseURL string
var testMongoStruct *test_funcs.TestMongoStruct

var userId = 12345
var userEmail = "testuser@example.com"

func TestMain(m *testing.M) {

	godotenv.Load(".env.test")

	testMongoStruct = test_funcs.SetUpMongoTestDatabase()
	defer testMongoStruct.Disconnect()

	// 起動前のセットアップ
	gin.SetMode(gin.TestMode)
	r := SetupRouter()
	testServer = &http.Server{
		Addr:    ":8881",
		Handler: r,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// サーバが立ち上がるまで待つ（またはヘルスチェックする）
	time.Sleep(200 * time.Millisecond)
	baseURL = "http://localhost:8881"

	fmt.Println("Test server started at", baseURL)

	// 全テスト実行
	exitCode := m.Run()

	// サーバをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = testServer.Shutdown(ctx)

	os.Exit(exitCode)
}

func createJwt(t *testing.T) string {
	jwt_secret := os.Getenv("JWT_SECRET")

	jwt, err := test_funcs.CreateMockJwtToken(
		userId,
		userEmail,
		time.Now().Add(1*time.Hour),
		[]byte(jwt_secret),
	)

	assert.NoError(t, err)

	return jwt
}

func createCsrf() string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := test_funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func request(method string, url string, body io.Reader, t *testing.T) (*http.Response, func() error) {
	jwt := createJwt(t)
	csrf := createCsrf()

	client := &http.Client{}
	requestUrl := baseURL + url
	fmt.Println("Request URL:", requestUrl)
	req, err := http.NewRequest(method, requestUrl, body)
	if method != "GET" && err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwt)
	if method != "GET" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, resp.Body.Close
}

func TestHealth(t *testing.T) {
	resp, close := request("GET", "/health", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "healthy")
}

func TestRoomCreate(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	body := strings.NewReader(`{"name":"TestRoom","is_private":false}`)
	resp, close := request("POST", "/room_create", body, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "Room created successfully")

	filter := bson.M{
		"name":    "TestRoom",
		"ownerid": userId,
		"members": userId,
	}

	exists, err := testMongoStruct.ExistContents(model.RoomCollectionName, filter)
	assert.NoError(t, err)
	assert.True(t, exists)

	count, err := testMongoStruct.CountContents(model.RoomCollectionName, filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestRoomJoin(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成
	createRoom := bson.M{
		"name":    "JoinableRoom",
		"ownerid": 99999,
		"members": []int{
			99999,
		},
	}
	room, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createRoom)
	assert.NoError(t, err)

	roomId := room.InsertedID.(primitive.ObjectID).Hex()
	body := strings.NewReader(`{"room_id":"` + roomId + `"}`)
	resp, close := request("POST", "/room_join", body, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "Joined room successfully")

	// ルームにユーザが追加されていることを確認
	filter := bson.M{
		"_id": room.InsertedID,
		"members": bson.M{
			"$in": []int{
				99999,
				userId,
			},
		},
	}
	exists, err := testMongoStruct.ExistContents(model.RoomCollectionName, filter)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRoomList(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成(複数パターン)
	// １．プライベートルーム自分がオーナーのルーム(表示される)
	// ２．プライベートルーム自分がメンバーのルーム(表示される)
	// ３．プライベートルーム（自分は関係ない:表示されない）
	// 4.パブリックルーム（自分は関係ない:表示される）
	variations := []model.Room{
		{Name: "PrivateRoom_Owner", OwnerID: userId, IsPrivate: true, Members: []int{userId}},
		{Name: "PrivateRoom_Member", OwnerID: 99999, IsPrivate: true, Members: []int{99999, userId}},
		{Name: "PrivateRoom_None", OwnerID: 88888, IsPrivate: true, Members: []int{88888}},
		{Name: "PublicRoom_None", OwnerID: 77777, IsPrivate: false, Members: []int{77777}},
	}

	for _, room := range variations {
		_, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, room)
		assert.NoError(t, err)
	}

	resp, close := request("GET", "/room_list", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)

	// 表示されるルームを確認
	assert.Contains(t, bodyString, "PrivateRoom_Owner")
	assert.Contains(t, bodyString, "PrivateRoom_Member")
	assert.Contains(t, bodyString, "PublicRoom_None")
	// 表示されないルームを確認
	assert.NotContains(t, bodyString, "PrivateRoom_None")
}

func TestPostChatMessage(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成
	createRoom := bson.M{
		"name":    "ChatRoom",
		"ownerid": userId,
		"members": []int{
			userId,
		},
	}
	room, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createRoom)
	assert.NoError(t, err)

	roomId := room.InsertedID.(primitive.ObjectID).Hex()
	body := strings.NewReader(`{"room_id":"` + roomId + `","message":"Hello, this is a test message."}`)
	resp, close := request("POST", "/post_chat_message", body, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "Chat posted successfully")

	// メッセージが保存されていることを確認
	filter := bson.M{
		"roomid":  roomId,
		"userid":  userId,
		"message": "Hello, this is a test message.",
	}

	exists, err := testMongoStruct.ExistContents(model.ChatMessageCollectionName, filter)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestLoadChat(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成
	createRoom := model.Room{
		Name:      "LoadChatRoom",
		OwnerID:   userId,
		IsPrivate: false,
		Members:   []int{userId},
	}
	room, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createRoom)
	assert.NoError(t, err)
	// チャットメッセージを作成
	chatMessages := []model.ChatMessage{
		{RoomID: room.InsertedID.(primitive.ObjectID).Hex(), UserID: userId, Message: "First test message", CreatedAt: time.Now()},
		{RoomID: room.InsertedID.(primitive.ObjectID).Hex(), UserID: userId, Message: "Second test message", CreatedAt: time.Now()},
	}

	for _, chat := range chatMessages {
		_, err := testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, chat)
		assert.NoError(t, err)
	}

	createNoiseRoom := model.Room{
		Name:      "NoiseRoom",
		OwnerID:   99999,
		IsPrivate: true,
		Members:   []int{99999, userId},
	}
	noiseRoom, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createNoiseRoom)
	assert.NoError(t, err)

	// ノイズとなるチャットメッセージを作成
	noiseChat := model.ChatMessage{
		RoomID:    noiseRoom.InsertedID.(primitive.ObjectID).Hex(),
		UserID:    99999,
		Message:   "Noise message",
		CreatedAt: time.Now(),
	}
	_, err = testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, noiseChat)
	assert.NoError(t, err)

	roomId := room.InsertedID.(primitive.ObjectID).Hex()
	resp, close := request("GET", "/load_chat/"+roomId, nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	fmt.Println("LoadChat Response Body:", bodyString)
	assert.Contains(t, bodyString, "First test message")
	assert.Contains(t, bodyString, "Second test message")
	assert.NotContains(t, bodyString, "Noise message")
}

func TestReadChatMessage(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成
	createRoom := model.Room{
		Name:      "ReadChatRoom",
		OwnerID:   userId,
		IsPrivate: false,
		Members:   []int{userId},
	}
	room, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createRoom)
	assert.NoError(t, err)

	messageIds := []primitive.ObjectID{}
	for i := 1; i <= 3; i++ {
		createChat := model.ChatMessage{
			RoomID:        room.InsertedID.(primitive.ObjectID).Hex(),
			UserID:        12345,
			Message:       fmt.Sprintf("Message %d", i),
			CreatedAt:     time.Now(),
			IsReadUserIds: []int{99999},
		}
		msg, err := testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, createChat)
		assert.NoError(t, err)
		messageIds = append(messageIds, msg.InsertedID.(primitive.ObjectID))
	}
	// ノイズメッセージを追加
	createNoiseChat := model.ChatMessage{
		RoomID:        room.InsertedID.(primitive.ObjectID).Hex(),
		UserID:        12345,
		Message:       "Noise message",
		CreatedAt:     time.Now(),
		IsReadUserIds: []int{99999},
	}
	noiseChat, err := testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, createNoiseChat)
	assert.NoError(t, err)

	chatIDListJson := []string{}
	for _, id := range messageIds {
		chatIDListJson = append(chatIDListJson, `"`+id.Hex()+`"`)
	}
	chatIDListStr := "[" + strings.Join(chatIDListJson, ",") + "]"

	roomId := room.InsertedID.(primitive.ObjectID).Hex()
	body := strings.NewReader(`{"room_id":"` + roomId + `","chat_id_list":` + chatIDListStr + `}`)
	resp, close := request("POST", "/read_chat", body, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "Chat messages marked as read")

	// メッセージが既読になっていることを確認
	for _, msgId := range messageIds {
		exist, err := testMongoStruct.ExistContents(model.ChatMessageCollectionName, bson.M{
			"_id": msgId,
			"isreaduserids": bson.M{
				"$in": []int{
					99999,
					userId,
				},
			},
		})
		assert.NoError(t, err)
		assert.True(t, exist)
	}

	// ノイズメッセージが既読になっていないことを確認
	exist, err := testMongoStruct.ExistContents(model.ChatMessageCollectionName, bson.M{
		"_id": noiseChat.InsertedID,
		"isreaduserids": bson.M{
			"$in": []int{
				userId,
			},
		},
	})
	assert.NoError(t, err)
	assert.False(t, exist)
}

func TestDeleteChatMessage(t *testing.T) {
	testMongoStruct.MongoCleanUp()

	// 事前にルームを作成
	createRoom := model.Room{
		Name:      "DeleteChatRoom",
		OwnerID:   userId,
		IsPrivate: false,
		Members:   []int{userId},
	}
	room, err := testMongoStruct.DB.Collection(model.RoomCollectionName).InsertOne(testMongoStruct.Ctx, createRoom)
	assert.NoError(t, err)

	// チャットメッセージを作成
	createChat := model.ChatMessage{
		RoomID:    room.InsertedID.(primitive.ObjectID).Hex(),
		UserID:    userId,
		Message:   "Message to be deleted",
		CreatedAt: time.Now(),
	}
	chat, err := testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, createChat)
	assert.NoError(t, err)
	// ノイズ用
	noiseCreateChat := model.ChatMessage{
		RoomID:    room.InsertedID.(primitive.ObjectID).Hex(),
		UserID:    userId,
		Message:   "Noise message",
		CreatedAt: time.Now(),
	}
	noiseChat, err := testMongoStruct.DB.Collection(model.ChatMessageCollectionName).InsertOne(testMongoStruct.Ctx, noiseCreateChat)
	assert.NoError(t, err)

	// 削除は生成が間違っていると、そもそも通ってしまうので、一度ここで、存在確認を行う
	exist, err := testMongoStruct.ExistContents(model.ChatMessageCollectionName, bson.M{
		"_id": chat.InsertedID,
	})
	assert.NoError(t, err)
	assert.True(t, exist)

	chatId := chat.InsertedID.(primitive.ObjectID).Hex()
	roomId := room.InsertedID.(primitive.ObjectID).Hex()
	body := strings.NewReader(`{"room_id":"` + roomId + `","message_id":"` + chatId + `"}`)
	resp, close := request("DELETE", "/delete_chat_message", body, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "Message deleted successfully")

	// メッセージが削除されていることを確認
	exist, err = testMongoStruct.ExistContents(model.ChatMessageCollectionName, bson.M{
		"_id": chat.InsertedID,
	})
	assert.NoError(t, err)
	assert.False(t, exist)

	// ノイズメッセージが削除されていないことを確認
	exist, err = testMongoStruct.ExistContents(model.ChatMessageCollectionName, bson.M{
		"_id": noiseChat.InsertedID,
	})
	assert.NoError(t, err)
	assert.True(t, exist)
}
