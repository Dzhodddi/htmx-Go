package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"project/internal/store"
)

var usernames = []string{
	"alphaWolf", "blueTiger", "cosmicRay", "deltaBlade", "echoStorm",
	"frozenPeak", "goldenFalcon", "hyperNova", "ironFist", "jadeShadow",
	"kingCobra", "lunarEcho", "midnightRider", "novaSpark", "omegaKnight",
	"pixelDrift", "quantumLeap", "redVortex", "silentArrow", "turboCharger",
	"ultraSonic", "venomStrike", "wildFlame", "xenoGhost", "youngViking",
	"zephyrWind", "amberGale", "binaryStar", "cyberWolf", "darkHorizon",
	"emberDash", "frostByte", "glitchGuru", "hazardZone", "icePhoenix",
	"jokerByte", "kryptoKnight", "laserHawk", "matrixWanderer", "neonGlider",
	"orbitRush", "phantomFlame", "quantumBolt", "ravenEye", "shadowPulse",
	"thunderWhale", "umbraKnight", "vortexRider", "warpSignal", "zetaCore",
}

var titles = []string{"5 Ways to Boost Productivity",
	"The Future of Remote Work",
	"Simple Tips for Better Sleep",
	"Understanding Cloud Storage",
	"How to Start a Side Hustle",
	"Mastering Time Management",
	"Beginner's Guide to Investing",
	"Top 10 Coding Resources",
	"Design Trends in 2025",
	"Healthy Eating on a Budget",
	"The Power of Daily Routines",
	"Why Mindfulness Matters",
	"Creating a Minimalist Workspace",
	"How to Learn Anything Faster",
	"Staying Motivated Long-Term",
	"Balancing Work and Life",
	"Intro to Web Development",
	"Tips for Better Focus",
	"The Art of Saying No",
	"Tools Every Freelancer Needs",
}

var contents = []string{
	"Actionable Tip: Elevate Your [Niche] with This Immediate Implement!",
	"Essential Resource: My Current Favorite [Tool/Resource] - Here's Why It's Key.",
	"Avoid This Pitfall: A Common [Niche] Mistake and Its Simple Resolution.",
	"Insightful Reflection: A Quote That Resonates & Its Relevance to You.",
	"Behind the Scenes: A Brief Look into My Process/Work/Recent Activity.",
	"Quick How-To: [Simple Task] in [Number] Concise Steps.",
	"Industry Snapshot: Key News/Trends & My Brief Perspective.",
	"Community Engagement: Your Thoughts on This [Niche] Question?",
	"Mini Case Study: A Short Success Story/Learning Moment.",
	"Audience Poll: Quick Question - Share Your Perspective!",
}

var commentsSlice = []string{
	"// A short, actionable tip related to the blog's specific topic.",
	"// A recommendation for a useful tool, website, or resource.",
	"// Highlights a common error within the niche and how to prevent it.",
	"// Shares an inspiring or relevant quote with a brief explanation.",
	"// Offers a glimpse into the author's process or behind-the-scenes activities.",
	"// A concise guide on how to do something specific in a few steps.",
	"// A brief summary of important news or trends in the industry.",
	"// Poses a question to encourage reader interaction and discussion.",
	"// A short story or example illustrating a key concept or success.",
	"// A simple poll to gather audience opinions on a specific topic.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, i := range users {
		if err := store.Users.Create(ctx, tx, i); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user", err)

		}
	}
	tx.Commit()
	posts := generatePosts(200, users)
	for _, i := range posts {
		if err := store.Posts.Create(ctx, i); err != nil {
			log.Println("Error creating post", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, i := range comments {
		if err := store.Comments.Create(ctx, i); err != nil {
			log.Println("Error creating comment", err)
			return
		}
	}
	log.Println("Seeding is completed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				titles[rand.Intn(len(titles))],
				titles[rand.Intn(len(titles))],
				titles[rand.Intn(len(titles))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: commentsSlice[rand.Intn(len(commentsSlice))],
		}
	}
	return comments
}
