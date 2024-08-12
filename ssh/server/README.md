# Notes from "Go SSH server complete example"


## `main()`

1. Create a `ssh.ServerConfig`
  - with the authn callback
  - read the server private key, and add it to the ServerConfig with `AddHostKey`

2. Create a `net.Listener`

3. In a for loop, `Accept` connections from the listener.
  - Using the connection from Accept, and the ServerConfig, create an ssh server connection using `ssh.NewServerConn()`.
  - Start a go routine to service the Request channel. This go routine discards the requests uing `ssh.DiscardRequests()`.
  - Start a go routine to service the NewChannel channel. See `handleChannel` below.

## `handleChannel()`

1. If the channel type of the `NewChannel` isn't "session", it `Reject`s the `NewChannel` and returns.

2. It `Accept`s the `NewChannel`. This returns a `Channel` and a channel `<- chan *ssh.Request`, or an error.

3. A bash command is created using `exec.Command("bash")`, and started using a pty, We then connect the 
   `Channel`'s Reader and Writer to the pty's `*os.File` Reader and Writer.

4. A go routine is created to service the request channel, handling
  - "shell", only `Reply`ing true to payload that don't contain a command
  - "pty-req", reseting the window size
  - "window-change", resting the window size



# Notes from "https://pkg.go.dev/golang.org/x/crypto/ssh"

## [ServerConfig](https://pkg.go.dev/golang.org/x/crypto/ssh#ServerConfig)

The immediately interesting fields of this structure are the authenticate callback functions that can be set.

- 	`PublicKeyCallback func(conn ConnMetadata, key PublicKey) (*Permissions, error)`
- 	`PasswordCallback func(conn ConnMetadata, password []byte) (*Permissions, error)`
- 	`KeyboardInteractiveCallback func(conn ConnMetadata, client KeyboardInteractiveChallenge) (*Permissions, error)`

Most interesting most likely being the `PublicKeyCallback`, since the idea here is to store the client's
public key is a database and use this callback function to verify the key. See the 
[NewServerConn example](https://pkg.go.dev/golang.org/x/crypto/ssh#example-NewServerConn) for a good
use of the `PublicKeyCallback`.

**Exercise:** Create a "hello world" ssh server and authenticate with a public key.

Also of interest:
-   `MaxAuthTries int`                    // If set to zero, the number of attempts are limited to 6.

-   `NoClientAuth boolNoClientAuth bool`  // NoClientAuth is true if clients are allowed to connect without authenticating.
-   `NoClientAuthCallback func(ConnMetadata) (*Permissions, error)`


### [ServerConfig.AddHostKey](https://pkg.go.dev/golang.org/x/crypto/ssh#ServerConfig.AddHostKey)
AddHostKey adds a private key as a host key. If an existing host key exists with the same public key format, it is replaced. Each server config must have at least one host key.

This is required for a server. See `Signer`/`ParsePrivateKey` below. These are used to parse the key from a file, which is then fed to AddHostKey.

## [Signer](https://pkg.go.dev/golang.org/x/crypto/ssh#Signer)

### [ParsePrivateKey](https://pkg.go.dev/golang.org/x/crypto/ssh#ParsePrivateKey)
Used to parse a private key from a series of bytes. Can be used when setting up a server's private host key config.


## [ServerConn](https://pkg.go.dev/golang.org/x/crypto/ssh#ServerConn)

### [NewServerConn](https://pkg.go.dev/golang.org/x/crypto/ssh#NewServerConn)
This returns a `*ssh.ServerConn`, a channel `<- chan ssh.NewChannel`, and another channel `<- chan *ssh.Request`.

The underlying transport is the tcp connection. 

The channel of `NewChannel` delivers new ssh channels being created.

The channel of `Request` delivers out of band requests.


## [NewChannel](https://pkg.go.dev/golang.org/x/crypto/ssh#NewChannel)
NewChannel represents an incoming request to a channel. It must either be accepted for use by calling `Accept`, or rejected by calling `Reject`.

## [Channel](https://pkg.go.dev/golang.org/x/crypto/ssh#Channel)
Channel is an interface that satisfies `io.ReadWriteCloser`. It has three additional methods.


# Exercises
- [ ] Build and run the gist from [Go SSH server complete example](https://gist.github.com/jpillora/b480fde82bff51a06238)
- [ ] Create a "hello world" ssh server and authenticate with a public key.


# TODO
- Read: https://blog.gopheracademy.com/go-and-ssh/
- https://github.com/jpillora/sshd-lite
- https://pkg.go.dev/github.com/paastech-cloud/git-ssh-server@v1.0.1

# References

- [Go SSH server complete example](https://gist.github.com/jpillora/b480fde82bff51a06238), - Read more here https://blog.gopheracademy.com/go-and-ssh/

- https://pkg.go.dev/golang.org/x/crypto/ssh
 - PROTOCOL: https://cvsweb.openbsd.org/cgi-bin/cvsweb/src/usr.bin/ssh/PROTOCOL?rev=HEAD
 - PROTOCOL.certkeys: http://cvsweb.openbsd.org/cgi-bin/cvsweb/src/usr.bin/ssh/PROTOCOL.certkeys?rev=HEAD
 - SSH-PARAMETERS:    http://www.iana.org/assignments/ssh-parameters/ssh-parameters.xml#ssh-parameters-1
