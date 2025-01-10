package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/munene-m/pepa/internal/models"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
	SessionStore *sessions.CookieStore
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	// Use a secure random key for session encryption
	authKey := []byte(os.Getenv("SESSION_KEY"))
	if len(authKey) == 0 {
		log.Fatal("SESSION_KEY environment variable is not set")
	}

	return &AuthHandler{
		DB:            db,
		SessionStore:  sessions.NewCookieStore(authKey),
	}
}

func InitializeGoogleAuth() {
	// Get credentials from environment variables
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY must be set")
	}

	// Initialize the session store for Gothic
	store := sessions.NewCookieStore([]byte(sessionKey))
	gothic.Store = store

	// Initialize Google provider
	googleProvider := google.New(googleClientID, googleClientSecret, callbackURL)
	
	// Register the provider with Goth
	goth.UseProviders(googleProvider)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := h.SessionStore.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Error getting session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear the session
	session.Values["user_id"] = nil
	session.Values["authenticated"] = false
	session.Values["provider"] = nil

	// Save the cleared session
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Error clearing session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to home or login page
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// func (h *AuthHandler)

func (h *AuthHandler) SignInWithProvider(w http.ResponseWriter, r *http.Request) {
    // provider := r.URL.Query().Get("provider")	
    
    // Add provider to query parameters
    q := r.URL.Query()
    q.Set("provider", "google")
    r.URL.RawQuery = q.Encode()

    // Use gothic to begin the auth handler
    gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
    user, err := gothic.CompleteUserAuth(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Extract name parts from email if FirstName/LastName are empty
    firstName := user.FirstName
    lastName := user.LastName
    if firstName == "" && lastName == "" && user.Email != "" {
        parts := strings.Split(user.Email, "@")
        nameParts := strings.Split(parts[0], ".")
        if len(nameParts) > 1 {
            firstName = strings.Title(nameParts[0])
            lastName = strings.Title(nameParts[1])
        } else {
            firstName = strings.Title(nameParts[0])
        }
    }

    // Create or update user in database with all available fields
    dbUser := &models.User{
        Email:          user.Email,
        Name:          user.Name,
        FirstName:     firstName,
        LastName:      lastName,
        GoogleID:      user.UserID,
        AccessToken:   user.AccessToken,
        RefreshToken:  user.RefreshToken,
        ProfilePicture: user.AvatarURL,
        Locale:         user.Location,
        EmailVerified:  true,
    }

    // Use Upsert to update all fields if user exists
    result := h.DB.Where("google_id = ?", user.UserID).
        Assign(dbUser).
        FirstOrCreate(dbUser)

    if result.Error != nil {
        // log.Printf("Database error: %v", result.Error)
        http.Error(w, "Error creating/updating user", http.StatusInternalServerError)
        return
    }

    var savedUser models.User
    h.DB.First(&savedUser, dbUser.ID)

    // Set session with additional data
    session, _ := h.SessionStore.Get(r, "auth-session")
    session.Values["user_id"] = dbUser.ID
    session.Values["authenticated"] = true
    session.Values["access_token"] = user.AccessToken
    session.Values["expires_at"] = user.ExpiresAt.Unix()
    session.Options.MaxAge = 86400 * 7 // 7 days
    if err := session.Save(r, w); err != nil {
        log.Printf("Session save error: %v", err)
    }
    http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := h.SessionStore.Get(r, "auth-session")
        auth, ok := session.Values["authenticated"].(bool)
        
        if !ok || !auth {
            http.Redirect(w, r, "/auth/google/signin", http.StatusTemporaryRedirect)
            return
        }
        
        next.ServeHTTP(w, r)
    }
}

func (h *AuthHandler) Success(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf(`
      <div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfull signed in!</h1>
          
          </div>
      </div>
  `)))
}