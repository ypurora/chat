
## How to Run


### Running Containers Manually

1. **Create a Docker network:**

    ```sh
    docker network create chatnetwork
    ```

2. **Build the server container:**

    ```sh
    docker build -t chatserver -f server/Dockerfile.server .
    ```

3. **Run the server container:**

    ```sh
    docker run -p 8080:8080 --network chatnetwork --name chatserver chatserver
    ```

4. **Build the client container:**

    ```sh
    docker build -t chatclient -f client/Dockerfile.client .
    ```

5. **Run the client container:**

    ```sh
    docker run -it --network chatnetwork -e SERVER_ADDR=chatserver:8080 chatclient
    ```

## Usage

### Register

1. When prompted, enter your desired username.
2. Enter your desired password.
3. Choose to register by typing `r` and pressing Enter.

### Login

1. When prompted, enter your username.
2. Enter your password.
3. Choose to log in by typing `l` and pressing Enter.

### Chat Commands

- **Create a Room**: Use the command `/create room_name` to create a new chat room.
- **Join a Room**: Use the command `/join room_name` to join an existing chat room.
- **Leave a Room**: Use the command `/leave room_name` to leave a chat room.
- **Direct Message**: Use the format `@username message` to send a direct message to a user.
- **Send Message**: Simply type your message and press Enter to send it to the current chat room.

### Chat Storage

- **Chat** is being stored in `output/{room_name}`
- **DM** is being stored in `output/dm/{user1_user2}`

### Example

```plaintext
Welcome to the Chat App
Enter username: alice
Enter password: secret
Do you want to (r)egister or (l)ogin? r
Successfully registered
Type messages to send to the chat. Type 'exit' to quit.
To send a DM, use the format: @username your message
To create a room, use the command: /create room_name
To join a room, use the command: /join room_name
To leave a room, use the command: /leave room_name

/create myroom
Joined room: myroom
@bob Hello, Bob!
bob: Hello, Alice!
```
