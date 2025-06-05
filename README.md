# Description

Facebook-like Social Network
A full-featured social networking platform built with Go (backend) and a modern JavaScript framework (frontend). This project includes features like posts, profiles, followers, private messaging, group chats, and notifications—all containerized using Docker.


# 🚀 Features
🔒 Authentication

Registration & login using sessions and cookies

User profiles: public/private toggle, optional avatar/nickname/about

👤 Followers

Follow/unfollow with request/approval for private profiles

🧾 Posts & Comments

Public, followers-only, and selected-followers post privacy

Image/GIF support

👥 Groups

Create, invite, request to join

Group posts/comments visible only to members

Event creation with (Going / Not Going)

💬 Chats

Private 1:1 messages

Group chat for group members

Emoji support

🔔 Notifications

Follower requests, group invites, join requests, new events, chat messages

🐳 Dockerized

Separate containers for frontend and backend

# 🧱 Tech Stack
Backend
Language: Go

Database: SQLite with golang-migrate

WebSocket: gorilla/websocket

Auth: Sessions & Cookies

Image Handling: JPEG, PNG, GIF

Frontend
Languages: JavaScript, HTML, CSS

Framework: Next.js

# 🐳 Docker Setup
Each service runs in its own container:

frontend: serves client-side application (via port 3000)

backend: serves API and WebSocket (via port 8080)

# Frontend
This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).


## Running the app with docker
To run the app using docker
1. Make  the startup script executable:

```bash
chmod +x start.sh
```
2. Start the app:
```bash
./start.sh
```
## Running the app without docker
To run the frontend locally without Docker:

1. Navigate to the frontend directory:
```bash
cd frontend
```
2. Install dependencies:
```bash
npm install
```
3. Start the development server:
```bash
npm run dev
```
Or, if you're using a different package manager:
```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.


This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!


# Backend

open another terminal
```bash
cd backend
go run .
```