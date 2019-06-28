# Pharmeum-api

### Environment variables
An example of environment variables

```dockerfile
# HTTP serve settings 
ENV 

# configuration of log level
ENV API_LOG_LEVEL=info
API_HTTP_ADDRESS=http://0.0.0.0:8080
# Database configuration
ENV APT_DB_URL=postgresql://postgres:1234567@db:5432/auth_service_db?sslmode=disable

#SMTP client configuration
ENV API_EMAIL_SENDER_ADDRESS=support@pharmeum.com
ENV API_EMAIL_SENDER_PASSWORD=12345678 
ENV API_SMTP_SERVER_HOST=google.com
ENV API_SMTP_SERVER_PORT=54321
```



[![N|Solid](https://s.dou.ua/CACHE/images/img/static/companies/HighChain_LOGO/57b91e978bda133de1e4148a5dae8e2b.png)](http://highchain.io/)

### Request and response

Api have four queries:
* Query which register new user.
    Method checks if existing user in database. If user not exist, method make register operation and sent confirm message on user email.
    If user already exist, method return code 400.
    
    Use link to make query: **POST (http://localhost:8080/user/signup)**
    Example body of query:
    
    ```json
    {
      "credentials": 
        {
          "email": "email@example.com",
           "password": "password"
        },
      "name": "Bob",
      "phone": "+1234567890",
      "country": "US",
      "city": "New York"
    }
    ```

* Query which make user active 
   Use link to make query: **GET (http://localhost:8080/user/confirm_email)**
   Example of query:
   http://localhost:8080/user/confirm_email?token=d78f8820-2e9e-4e1a-a5a3-46c1a0e5c400

* Query which execute changing password and sent confirm mail on user email
    Use link to make query: GET (http://localhost:8080/user/confirm_email)
    Example body of query:
    ```json
    {
      "password": "1234567890"
    }
    ```
    
* **Reset password**
    This method generate token and sent recovery link on your mail
    Example body of query:    
    ```json
    {
      "email": "email@example.com"
    }
    ```
    Use link to make query: **GET (http://localhost:8080/user/confirm_email)**
