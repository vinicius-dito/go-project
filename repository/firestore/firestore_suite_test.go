package user_firestore_repository_test

import (
	"context"
	"go-project/domain"
	user_firestore_repository "go-project/repository/firestore"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func setup() (context.Context, *firestore.CollectionRef) {
	ctx := context.Background()

	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	firestoreClient, err := firestore.NewClient(ctx, "demo-teste")
	if err != nil {
		panic(err)
	}

	firestoreUserCollection := firestoreClient.Collection("teste")

	return ctx, firestoreUserCollection
}

func TestUserFirestoreRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firestore Repository Suite")
}

var _ = Describe("UserFirestoreRepository", func() {
	var userFirestoreRepository user_firestore_repository.UserFirestoreRepositoy
	userID := "123"
	ctx, firestoreUserCollection := setup()

	It("should return error when missing dependencies", func() {
		_, err := user_firestore_repository.NewUsersFirestoreRepository(nil)

		Expect(err).Should(HaveOccurred())
		Expect(err).To(MatchError("failed to create firestore repository"))
	})

	It("should succeed to create a new repository when all dependencies are correct", func() {
		_, err := user_firestore_repository.NewUsersFirestoreRepository(firestoreUserCollection)

		Expect(err).ShouldNot(HaveOccurred())
	})

	BeforeEach(func() {
		userFirestoreRepository, _ = user_firestore_repository.NewUsersFirestoreRepository(firestoreUserCollection)
	})

	AfterEach(func() {
		// fazer uma função pra limpar o banco
		//Delete user
		firestoreUserCollection.Doc(userID).Delete(ctx)
	})

	It("should succeed to save user", func() {
		//Act
		err := userFirestoreRepository.Save(ctx, domain.User{UserId: userID})

		//Assert
		user, _ := userFirestoreRepository.Get(ctx, userID)

		Expect(err).ToNot(HaveOccurred())
		Expect(user.UserId).To(Equal(userID))
	})

	It("should fail to save user", func() {
		//Act
		err := userFirestoreRepository.Save(ctx, domain.User{UserId: ""})

		//Assert
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to insert user into Firestore"))
	})

	It("should succeed to get user", func() {
		//Arrange
		userFirestoreRepository.Save(ctx, domain.User{UserId: userID})

		//Act
		_, err := userFirestoreRepository.Get(ctx, userID)

		//Assert
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail to get user", func() {
		//Act
		_, err := userFirestoreRepository.Get(ctx, "")

		//Assert
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to get user from Firestore"))
	})

	It("should fail to get user if doc struct is wrong", func() {
		//Arrange
		type invalidUser struct {
			UserID int `firestore:"user_id"`
		}
		firestoreUserCollection.Doc(userID).Set(ctx, invalidUser{UserID: 123})

		//Act
		_, err := userFirestoreRepository.Get(ctx, userID)

		//Assert
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to parse Firestore document"))
	})
})
