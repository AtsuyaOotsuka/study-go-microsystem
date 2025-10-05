package chat_svc

import (
	"microservices/chat/internal/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConvertRoomList(t *testing.T) {
	rooms := []model.Room{
		{
			ID:        primitive.NewObjectID(),
			Name:      "is owner room",
			OwnerID:   1,
			IsPrivate: true,
			Members:   []int{1, 2, 3},
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "is member room",
			OwnerID:   2,
			IsPrivate: false,
			Members:   []int{1, 2, 3},
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "not member room",
			OwnerID:   3,
			IsPrivate: false,
			Members:   []int{2, 3},
			CreatedAt: time.Now(),
		},
	}

	svc := NewChatSvc()

	result := svc.ConvertRoomList(rooms, 1)

	if len(result) != 3 {
		t.Fatalf("expected 3 rooms, got %d", len(result))
	}

	expect := []struct {
		Name        string
		IsOwner     bool
		IsMember    bool
		MemberCount int
	}{
		{"is owner room", true, true, 3},
		{"is member room", false, true, 3},
		{"not member room", false, false, 2},
	}

	for i, r := range result {
		if r.Name != expect[i].Name {
			t.Errorf("expected Name %s, got %s", expect[i].Name, r.Name)
		}
		if r.IsOwner != expect[i].IsOwner {
			t.Errorf("expected IsOwner %v, got %v", expect[i].IsOwner, r.IsOwner)
		}
		if r.IsMember != expect[i].IsMember {
			t.Errorf("expected IsMember %v, got %v", expect[i].IsMember, r.IsMember)
		}
		if r.MemberCount != expect[i].MemberCount {
			t.Errorf("expected MemberCount %d, got %d", expect[i].MemberCount, r.MemberCount)
		}
	}
}
