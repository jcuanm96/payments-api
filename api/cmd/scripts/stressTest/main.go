package main

import (
	"fmt"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	vama "github.com/VamaSingapore/vama-api/internal/vamaClient"
	"github.com/sirupsen/logrus"
)

func main() {
	basePhoneNumber := "28179699"
	// createClientAndSpam(basePhoneNumber + "02")
	var wg sync.WaitGroup
	for i := 10; i < 99; i++ {
		phoneNumber := basePhoneNumber + fmt.Sprint(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			createClientAndSpam(phoneNumber)
		}()
	}

	wg.Wait()
}

func createClientAndSpam(phoneNumber string) {
	const baseURL = "127.0.0.1:8080"
	c := vama.NewClient(baseURL)

	myNumber := phoneNumber[len(phoneNumber)-2:]
	println("signing in user", myNumber)
	signUpOrSignIn(c, phoneNumber)

	_, followErr := c.Follow(request.Follow{UserID: 7})
	if followErr != nil {
		logrus.Errorf("User %s failed following: %s", myNumber, followErr.Error())
	}

	feedResp, feedErr := c.GetHomeFeedPosts(request.GetUserFeedPosts{})
	if feedErr != nil {
		logrus.Fatalf("User %s failed getting home feed posts: %s", myNumber, feedErr)
	}

	postID := feedResp.FeedPosts[0].ID
	_, getPostErr := c.GetFeedPostByID(postID)
	if getPostErr != nil {
		logrus.Fatalf("User %s failed getting feed post by ID: %s", myNumber, getPostErr)
	}

	makeCommentReq := request.MakeComment{
		Text:   fmt.Sprintf("User %s making a STATEMENT.", myNumber),
		PostID: postID,
	}
	makeCommentErr := c.MakeComment(makeCommentReq)
	if makeCommentErr != nil {
		logrus.Fatalf("User %s failed to make a comment: %s", makeCommentErr.Error())
	}

}

func signUpOrSignIn(c *vama.Client, phoneNumber string) {
	myNumber := phoneNumber[len(phoneNumber)-2:]
	signUpReq := request.SignupSMS{
		FirstName: "User" + myNumber,
		LastName:  "Name" + myNumber,
		Phone: request.Phone{
			CountryCode: "US",
			Number:      phoneNumber,
		},
		Code: "123456",
	}

	signUpResp, signUpErr := c.SignUpSMS(signUpReq)
	if signUpErr == nil {
		c.SetAccessToken(signUpResp.Credentials.AccessToken)
		return
	}

	signInReq := request.SignInSMS{
		Phone: signUpReq.Phone,
		Code:  "123456",
	}
	signInRes, signInErr := c.SignInSMS(signInReq)
	if signInErr != nil {
		logrus.Fatalf("Error signing in user %s: %s", myNumber, signInErr.Error())
		return
	}

	c.SetAccessToken(signInRes.Credentials.AccessToken)
}
