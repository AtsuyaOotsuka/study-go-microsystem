package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("⚠️ .env読み込み失敗:", err)
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func mongoConnect() (*mongo.Database, context.Context) {
	// タイムアウト付きのcontext
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// 接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("MongoDB 接続成功！")

	// データベースとコレクションの選択
	db := client.Database("chatapp")
	return db, ctx
}

func createRoom(title string) interface{} {
	db, ctx := mongoConnect()
	collection := db.Collection("rooms")

	room := map[string]interface{}{
		"title":      title,
		"created_at": time.Now(),
	}

	insertResult, err := collection.InsertOne(ctx, room)
	if err != nil {
		panic(err)
	}

	fmt.Printf("挿入成功！_id: %v\n", insertResult.InsertedID)

	return insertResult.InsertedID
}

func joinRoom(userID int, roomID interface{}) {
	db, ctx := mongoConnect()
	collection := db.Collection("user_rooms")

	userRoom := map[string]interface{}{
		"user_id": userID,
		"room_id": roomID,
	}

	insertResult, err := collection.InsertOne(ctx, userRoom)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ルーム参加成功！_id: %v\n", insertResult.InsertedID)
}

func getMessages(roomID interface{}) []map[string]interface{} {
	db, ctx := mongoConnect()
	collection := db.Collection("messages")

	filter := map[string]interface{}{"room_id": roomID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	var messages []map[string]interface{}
	for cursor.Next(ctx) {
		var message map[string]interface{}
		if err := cursor.Decode(&message); err != nil {
			panic(err)
		}
		messages = append(messages, message)
	}

	return messages
}

func sendMessage(roomID interface{}, userID int, message string) {
	db, ctx := mongoConnect()
	collection := db.Collection("messages")

	msg := map[string]interface{}{
		"room_id":        roomID,
		"sender_user_id": userID,
		"message":        message,
		"created_at":     time.Now(),
	}

	insertResult, err := collection.InsertOne(ctx, msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("メッセージ送信成功！_id: %v\n", insertResult.InsertedID)
}

func deleteMessage(messageID interface{}, userId int) {
	db, ctx := mongoConnect()
	collection := db.Collection("messages")

	filter := map[string]interface{}{"_id": messageID, "sender_user_id": userId}
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}

	if deleteResult.DeletedCount > 0 {
		fmt.Printf("メッセージ削除成功！_id: %v\n", messageID)
	} else {
		fmt.Println("メッセージが見つかりませんでした。")
	}
}

func convertIdToString(id string) interface{} {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("ID変換エラー:", err)
		return nil
	}
	return ObjectID
}

func main() {

	// ルーム生成
	http.HandleFunc("/create_room", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTでアクセスしてね！", http.StatusMethodNotAllowed)
			return
		}

		roomTitle := r.FormValue("title")
		if roomTitle == "" {
			http.Error(w, "ルームタイトルを入力してね！", http.StatusBadRequest)
			return
		}

		roomID := createRoom(roomTitle)
		if roomID == nil {
			http.Error(w, "ルーム作成に失敗しました！", http.StatusInternalServerError)
			return
		}
		joinRoom(r.Context().Value("userID").(int), roomID)

		sendMessage(roomID, r.Context().Value("userID").(int), "ルームが作成されました！")

		fmt.Fprintf(w, "ルーム作成成功！ルームID: %v", roomID)
	}))

	// ルームに参加
	http.HandleFunc("/join_room", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTでアクセスしてね！", http.StatusMethodNotAllowed)
			return
		}

		roomID := r.FormValue("room_id")
		if roomID == "" {
			http.Error(w, "ルームIDを入力してね！", http.StatusBadRequest)
			return
		}

		joinRoom(r.Context().Value("userID").(int), convertIdToString(roomID))

		sendMessage(convertIdToString(roomID), r.Context().Value("userID").(int), "ルームに参加しました！")
		fmt.Fprintf(w, "ルーム参加成功！ルームID: %s", roomID)
	}))

	// メッセージを確認
	http.HandleFunc("/get_messages", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GETでアクセスしてね！", http.StatusMethodNotAllowed)
			return
		}
		roomID := r.FormValue("room_id")
		if roomID == "" {
			http.Error(w, "ルームIDを入力してね！", http.StatusBadRequest)
			return
		}

		messages := getMessages(convertIdToString(roomID))
		if len(messages) == 0 {
			http.Error(w, "メッセージがありません！", http.StatusNotFound)
			return
		}
		for _, msg := range messages {
			fmt.Fprintf(w, "メッセージ: %s, 送信者ID: %d, 送信日時: %s\n",
				msg["message"], msg["sender_user_id"], msg["created_at"])
		}
		fmt.Fprintf(w, "合計 %d 件のメッセージがあります。", len(messages))
	}))

	// メッセージを送信
	http.HandleFunc("/send_message", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTでアクセスしてね！", http.StatusMethodNotAllowed)
			return
		}
		roomID := r.FormValue("room_id")
		if roomID == "" {
			http.Error(w, "ルームIDを入力してね！", http.StatusBadRequest)
			return
		}
		message := r.FormValue("message")
		if message == "" {
			http.Error(w, "メッセージを入力してね！", http.StatusBadRequest)
			return
		}
		sendMessage(convertIdToString(roomID), r.Context().Value("userID").(int), message)
		fmt.Fprintf(w, "メッセージ送信成功！ルームID: %s, メッセージ: %s", roomID, message)
	}))

	// メッセージを削除
	http.HandleFunc("/delete_message", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "DELETEでアクセスしてね！", http.StatusMethodNotAllowed)
			return
		}
		messageID := r.FormValue("message_id")
		if messageID == "" {
			http.Error(w, "メッセージIDを入力してね！", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value("userID").(int)
		deleteMessage(convertIdToString(messageID), userID)
		fmt.Fprintf(w, "メッセージ削除成功！メッセージID: %s, ユーザーID: %d", messageID, userID)
	}))

	http.ListenAndServe(":8081", nil)
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractBearerToken(r)
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 🎯 user_id を context にセット！
		idVal := claims["sub"].(float64)
		userID := int(idVal)
		ctx := context.WithValue(r.Context(), "userID", userID)
		// 🎯 email も context にセット
		email := claims["email"].(string)
		ctx = context.WithValue(ctx, "email", email)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func validateJWT(tokenString string) bool {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("Token parse error:", err)
		return false
	}

	// ここまで来たらほぼ成功
	fmt.Println("Parsed Claims:", claims)
	return token.Valid
}
