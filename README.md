# ShortenURL  
## Features  
- **Shorten URL**: Generate short, easy-to-share URLs with a unique identifier for long URLs.
- **Click Tracking**: Track the number of clicks and the last access time for each shortened URL.
- **Performance Boost**: Use Go's goroutines to handle redirects and analytics in parallel, improving performance and reducing latency when redirecting users to the original URL.
- **User Authentication**: Users can create accounts and log in to manage their shortened URLs, with authentication and session management.

## Demo  
- Demo video on YouTube: https://www.youtube.com/watch?v=BEi3t45F4Qs

## Tech Stack  
- **Go**: The core language used for building the application, known for its speed and efficiency.  
- **Goroutines**: Utilized for handling concurrent operations like redirects and analytics, optimizing performance.  
- **HTTP**: For handling web requests and responses.  
- **Mux**: A powerful HTTP router for Go, used for routing and managing URL patterns.  
- **GORM**: Object-relational mapping (ORM) tool for Go, used to interact with the PostgreSQL database.  
- **Template Rendering**: Used to generate dynamic HTML pages for the user interface.  
- **Sessions**: Implemented for user authentication and maintaining session state across requests.  
- **Unit Tests**: Ensuring the reliability of the application with unit testing for each component.  
- **PostgreSQL**: A relational database for storing URLs, user data, and tracking information.