package users_firestore_repository

import (
	"context"
	"fmt"
	"go-project/domain"

	"cloud.google.com/go/firestore"
)

type UsersFirestoreRepositoy struct {
	firestoreClient     firestore.Client
	firestoreCollection string
}

func NewUsersFirestoreRepository(firestoreClient firestore.Client, firestoreCollection string) UsersFirestoreRepositoy {
	return UsersFirestoreRepositoy{
		firestoreClient:     firestoreClient,
		firestoreCollection: firestoreCollection,
	}
}

func (ufr UsersFirestoreRepositoy) Save(ctx context.Context, user domain.Users) error {
	if _, err := ufr.firestoreClient.Collection(ufr.firestoreCollection).Doc(user.UserId).Set(ctx, user); err != nil {
		return fmt.Errorf("failed to insert user into Firestore: %v", err)
	}

	return nil
}

func (ufr UsersFirestoreRepositoy) Get(ctx context.Context, userID string) (domain.Users, error) {
	var user domain.Users

	userDoc, err := ufr.firestoreClient.Collection(ufr.firestoreCollection).Doc(userID).Get(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to get user from Firestore: %v", err)
	}

	if err = userDoc.DataTo(&user); err != nil {
		return user, fmt.Errorf("failed to parse Firestore document: %v", err)
	}

	return user, nil
}
