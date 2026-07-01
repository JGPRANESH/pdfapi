package services

import (
	"fmt"
	"math/rand"
	"time"
)

type NotificationTemplate struct {
	Title string
	Body  string
}

var quizTemplates = []NotificationTemplate{
	{
		Title: "🚨 New Mock Just Dropped!",
		Body:  "%s just landed. No excuses, go cook! 🔥",
	},
	{
		Title: "👀 Psst... New Challenge!",
		Body:  "%s is live. Let's see if you're built different. 😏",
	},
	{
		Title: "💀 Bro, It's Quiz Time",
		Body:  "Your new %s just dropped. Don't let it collect dust. 📚",
	},
	{
		Title: "⚡ Main Character Moment",
		Body:  "Crush the new %s mock and flex that score. 😎",
	},
	{
		Title: "🔥 Wake Up Bestie",
		Body:  "%s is waiting. Time to lock in. 🎯",
	},
	{
		Title: "🎮 XP Farming Time",
		Body:  "Complete the new %s mock and level up your skills. 🚀",
	},
	{
		Title: "📢 Certified Brain Rot Break",
		Body:  "Enough scrolling. Your %s quiz is live. 🧠✨",
	},
	{
		Title: "😤 You Got This",
		Body:  "Think you can ace %s? Prove it. 💯",
	},
	{
		Title: "🎯 One More Mock?",
		Body:  "%s just dropped. Future you will thank you. 🙌",
	},
	{
		Title: "🚀 Go Get That W",
		Body:  "A fresh %s is live. Time to secure the win. 🏆",
	},
	{
		Title: "🧠 Brain Buff Activated",
		Body:  "Jump into the new %s and stack that knowledge. 📈",
	},
	{
		Title: "💥 No Cap",
		Body:  "%s is live and kinda slaps. Check it out. 👀",
	},
	{
		Title: "✨ Your Sign to Study",
		Body:  "Yes, this is it. Start the new %s mock now. 📚",
	},
	{
		Title: "🎉 Fresh Drop Alert",
		Body:  "%s just went live. First attempt = best attempt? 🤔",
	},
	{
		Title: "🏁 Race Against Yourself",
		Body:  "Beat your last score in the new %s mock. Let's go! 🔥",
	},
}

func GetRandomQuizNotification(examName string) (string, string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	template := quizTemplates[r.Intn(len(quizTemplates))]

	return template.Title, fmt.Sprintf(template.Body, examName)
}
