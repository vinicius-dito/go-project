package user_firestore_repository

import (
	"context"
	"fmt"
	"go-project/domain"

	"cloud.google.com/go/firestore"
)

type UserFirestoreRepositoy struct {
	firestoreCollection *firestore.CollectionRef
}

func NewUsersFirestoreRepository(firestoreCollection *firestore.CollectionRef) UserFirestoreRepositoy {
	return UserFirestoreRepositoy{
		firestoreCollection: firestoreCollection,
	}
}

func (ufr UserFirestoreRepositoy) Save(ctx context.Context, user domain.User) error {
	if _, err := ufr.firestoreCollection.Doc(user.UserId).Set(ctx, user); err != nil {
		return fmt.Errorf("failed to insert user into Firestore: %v", err)
	}

	return nil
}

func (ufr UserFirestoreRepositoy) Get(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

	userDoc, err := ufr.firestoreCollection.Doc(userID).Get(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to get user from Firestore: %v", err)
	}

	if err = userDoc.DataTo(&user); err != nil {
		return user, fmt.Errorf("failed to parse Firestore document: %v", err)
	}

	return user, nil
}
