Kelimelik

Description

Kelimelik is a dynamic blogging platform designed for writers and enthusiasts who want to share their thoughts on any topic. Built with the Go programming language, Kelimelik offers a fast, lightweight, and scalable solution for creating, managing, and publishing blog content. Whether you're a casual blogger or a professional content creator, Kelimelik provides a user-friendly interface to express your ideas with ease.

Features





Flexible Blogging: Create and publish blog posts on any topic without restrictions.



Fast Performance: Powered by Go, ensuring quick load times and efficient handling of traffic.



Markdown Support: Write posts using Markdown for easy formatting.



Customizable Themes: Choose or create themes to personalize the look and feel of your blog.



SEO-Friendly: Built-in features to optimize your posts for search engines.



User Authentication: Secure login and user management for authors and admins.



Responsive Design: Accessible on desktops, tablets, and mobile devices.

Installation





Clone the Repository:

git clone https://github.com/sengka/kelimelik.git
cd kelimelik



Install Dependencies: Ensure you have Go installed (version 1.16 or higher). Then run:

go mod tidy





Run the Application: Start the server:

go run main.go

The blog platform will be available at http://localhost:8080.

Usage





Creating a Post: Log in to the admin panel, navigate to the "New Post" section, and write your content using Markdown.



Publishing: Save drafts or publish posts instantly to make them live.



Customization: Modify themes or templates in the templates/ directory to match your style.



Managing Users: Admins can manage user accounts and permissions via the admin dashboard.

Development





Tech Stack:





Backend: Go (Golang)



Frontend: HTML, CSS, 



Database: sqlite




