package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/saifwork/mock-service/internal/core/config"
	database "github.com/saifwork/mock-service/internal/core/mongo"
	"github.com/saifwork/mock-service/internal/dtos"
	"github.com/saifwork/mock-service/internal/models"
	"github.com/saifwork/mock-service/internal/utils"
)

type AuthService struct {
	coll *mongo.Collection
	ctx  context.Context
	cfg  *config.Config
}

// Constructor
func NewAuthService(client *mongo.Client, cfg *config.Config) *AuthService {
	collection := client.Database(cfg.MongoDBName).Collection(database.Collections.Users)
	return &AuthService{
		coll: collection,
		ctx:  context.Background(),
		cfg:  cfg,
	}
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var user models.User
	err = s.coll.FindOne(s.ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.PasswordHash = "" // hide sensitive field
	return &user, nil
}

// -------------------- Signup --------------------
func (s *AuthService) Signup(req *dtos.UserRegisterRequestDto) error {

	// Check if email already exists
	count, _ := s.coll.CountDocuments(s.ctx, bson.M{"email": req.Email})
	if count > 0 {
		return errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:           primitive.NewObjectID(),
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashed),
		IsUpgraded:   false,
		IsVerified:   false,
		IsActive:     false,

		VerificationToken:   primitive.NewObjectID().Hex(),
		VerificationExpires: time.Now().Add(24 * time.Hour),

		LastVerificationSentAt: time.Now(),
		VerificationEmailCount: 1, // first email sent

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.coll.InsertOne(s.ctx, user)
	if err != nil {
		return err
	}

	// Mock email send
	link := fmt.Sprintf("%s/auth/verify-email?token=%s", s.cfg.AppBaseURL, user.VerificationToken)
	log.Printf("[EMAIL MOCK] Verification link for %s: %s\n", user.Email, link)

	subject := "Verify your email - MockNode"
	html := utils.BuildEmailTemplate(
		"MockNode",
		"Welcome to MockNode 👋",
		fmt.Sprintf(`Hey %s,<br><br>
	We’re thrilled to have you on board! 🎉<br><br>

	<strong>MockNode</strong> is your personal developer playground — a space where you can easily create, test, and share mock APIs for your projects. 
	Whether you’re building your next big idea, testing integrations, or just learning how APIs work, MockNode makes it all effortless and fun.<br><br>

	Here’s what you can do right away:
	<ul style="text-align:left; margin:auto; width:80%%; color:#c9d1d9;">
		<li>⚡ Instantly create mock endpoints with zero setup.</li>
		<li>📦 Organize your projects and manage collections easily.</li>
		<li>🌍 Test your API responses directly in your apps.</li>
	</ul>

	You’re just one step away from exploring the full power of MockNode.<br><br>
	Click the button below to verify your email and start your journey! 🚀<br><br>

	Welcome to the MockNode family — where developers build, test, and learn together.<br><br>
	Happy coding! 💻<br>
	The MockNode Team
	`, user.FullName),
		"Verify Email",
		link,
	)

	if err := utils.SendEmail(s.cfg, user.Email, subject, html); err != nil {
		log.Printf("❌ Failed to send verification email: %v\n", err)
	} else {
		log.Printf("✅ Verification email sent to %s\n", user.Email)
	}

	return nil
}

// -------------------- Resend Verification Email --------------------
func (s *AuthService) ResendVerification(email string) error {

	// Find user by email
	var user models.User
	err := s.coll.FindOne(s.ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	// Already verified — no need to resend
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// ⛔ Cooldown 2 minutes
	if time.Since(user.LastVerificationSentAt) < 2*time.Minute {
		return errors.New("please wait before requesting another verification email")
	}

	// ⛔ Daily limit (reset if new day)
	if time.Now().Day() != user.LastVerificationSentAt.Day() {
		user.VerificationEmailCount = 0
	}

	if user.VerificationEmailCount >= 5 {
		return errors.New("daily limit reached — try again tomorrow")
	}

	// Generate new token + expiry
	newToken := primitive.NewObjectID().Hex()
	newExpiry := time.Now().Add(24 * time.Hour)

	// Update user with the new token
	_, err = s.coll.UpdateOne(
		s.ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"verificationToken":   newToken,
				"verificationExpires": newExpiry,
				"updatedAt":           time.Now(),
			},
		},
	)
	if err != nil {
		return err
	}

	// Create the verification link
	link := fmt.Sprintf("%s/auth/verify-email?token=%s", s.cfg.AppBaseURL, newToken)
	log.Printf("[EMAIL MOCK] (RESEND) Verification link for %s: %s\n", user.Email, link)

	subject := "Verify your email - MockNode"
	html := utils.BuildEmailTemplate(
		"MockNode",
		"Verify Your Email Again 👋",
		fmt.Sprintf(`Hey %s,<br><br>

It looks like you didn’t get a chance to verify your email earlier — no worries! 😊<br><br>

Click the button below to verify your email and unlock all features of <strong>MockNode</strong>.<br><br>

If you didn’t request this, you can safely ignore it.<br><br>
`, user.FullName),
		"Verify Email",
		link,
	)

	// Send email
	if err := utils.SendEmail(s.cfg, user.Email, subject, html); err != nil {
		log.Printf("❌ Failed to resend verification email: %v\n", err)
		return errors.New("failed to send verification email")
	}

	log.Printf("✅ Verification email re-sent to %s\n", user.Email)

	return nil
}

// -------------------- Verify Email --------------------
func (s *AuthService) VerifyEmail(token string) error {

	filter := bson.M{
		"verificationToken": token,
		"isVerified":        false,
		"isActive":          false,
		"verificationExpires": bson.M{
			"$gt": time.Now(),
		},
	}

	update := bson.M{
		"$set": bson.M{
			"isVerified": true,
			"isActive":   true,
			"updatedAt":  time.Now(),
		},
		"$unset": bson.M{
			"verificationToken":      "",
			"verificationExpires":    "",
			"lastVerificationSentAt": "",
			"verificationEmailCount": "",
		},
	}

	res, err := s.coll.UpdateOne(s.ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("invalid or expired token")
	}

	return nil
}

// -------------------- Login --------------------
func (s *AuthService) Login(req *dtos.UserLoginRequestDto) (string, error) {

	var user models.User
	err := s.coll.FindOne(s.ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), s.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

// -------------------- Forgot Password --------------------
func (s *AuthService) ForgotPassword(email string) error {

	token := primitive.NewObjectID().Hex()
	expiry := time.Now().Add(1 * time.Hour)

	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"resetToken":   token,
			"resetExpires": expiry,
			"updatedAt":    time.Now(),
		},
	}

	res, err := s.coll.UpdateOne(s.ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}

	link := fmt.Sprintf("%s/reset-forgot-password?token=%s", s.cfg.AppBaseURL, token)
	log.Printf("[EMAIL MOCK] Password reset link for %s: %s\n", email, link)

	subject := "Reset your password - MockNode"
	html := utils.BuildEmailTemplate(
		"MockNode",
		"Reset Your Password 🔒",
		"We received a request to reset your password. Click the button below to set a new one.",
		"Reset Password",
		link,
	)

	if err := utils.SendEmail(s.cfg, email, subject, html); err != nil {
		log.Printf("❌ Failed to send reset email: %v\n", err)
	} else {
		log.Printf("✅ Password reset email sent to %s\n", email)
	}

	return nil
}

// -------------------- Reset Password via Forgot --------------------
func (s *AuthService) ResetForgotPassword(token, newPassword string) error {

	var user models.User
	err := s.coll.FindOne(s.ctx, bson.M{
		"resetToken": token,
		"resetExpires": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&user)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	update := bson.M{
		"$set": bson.M{
			"passwordHash": string(hashed),
			"updatedAt":    time.Now(),
		},
		"$unset": bson.M{
			"resetToken":   "",
			"resetExpires": "",
		},
	}

	_, err = s.coll.UpdateByID(s.ctx, user.ID, update)
	return err
}

// -------------------- Reset Password (Logged-in user) --------------------
func (s *AuthService) ResetPassword(userID, oldPwd, newPwd string) error {

	objID, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	err := s.coll.FindOne(s.ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPwd)); err != nil {
		return errors.New("invalid old password")
	}

	// ✅ Check if new password is same as old
	if oldPwd == newPwd {
		return errors.New("new password cannot be same as old password")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)

	_, err = s.coll.UpdateByID(s.ctx, user.ID, bson.M{
		"$set": bson.M{
			"passwordHash": string(hashed),
			"updatedAt":    time.Now(),
		},
	})
	return err
}
