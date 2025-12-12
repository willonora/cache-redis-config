# cache-redis-config
=====================================

## Overview
cache-redis-config is a software project designed to provide a simple and efficient way to manage Redis cache configurations.

## Features
*   Supports multiple Redis instances
*   Provides a user-friendly interface for configuration management
*   Allows for easy switching between different Redis configurations

## Requirements
*   Node.js 16.0 or higher
*   Redis 6.0 or higher

## Installation
```bash
npm install
```

## Usage
```javascript
const RedisConfig = require('./src/RedisConfig');

const config = new RedisConfig({
    host: 'localhost',
    port: 6379,
});

config.connect().then(() => {
    console.log('Connected to Redis');
}).catch((err) => {
    console.error('Error connecting to Redis:', err);
});

// Set a value in Redis
config.set('key', 'value').then(() => {
    console.log('Value set successfully');
}).catch((err) => {
    console.error('Error setting value:', err);
});

// Get a value from Redis
config.get('key').then((value) => {
    console.log('Value:', value);
}).catch((err) => {
    console.error('Error getting value:', err);
});
```

## Configuration Options
The following configuration options are available:
*   `host`: The hostname or IP address of the Redis instance
*   `port`: The port number of the Redis instance
*   `username`: The username to use for authentication
*   `password`: The password to use for authentication
*   `db`: The Redis database to use

## Contributing
Contributions are welcome and encouraged. To contribute, please fork the repository and submit a pull request.

## License
cache-redis-config is licensed under the MIT License. See LICENSE for details.

## Changelog
### v1.0.0
*   Initial release
```javascript
// src/RedisConfig.js
class RedisConfig {
    constructor(options) {
        this.host = options.host;
        this.port = options.port;
        this.username = options.username;
        this.password = options.password;
        this.db = options.db;
    }

    connect() {
        return new Promise((resolve, reject) => {
            // Connect to Redis using the provided options
            const redis = require('redis');
            const client = redis.createClient({
                host: this.host,
                port: this.port,
                username: this.username,
                password: this.password,
                db: this.db,
            });

            client.on('connect', () => {
                resolve(client);
            });

            client.on('error', (err) => {
                reject(err);
            });
        });
    }

    set(key, value) {
        return new Promise((resolve, reject) => {
            this.connect().then((client) => {
                client.set(key, value, (err, reply) => {
                    if (err) {
                        reject(err);
                    } else {
                        resolve(reply);
                    }
                });
            }).catch((err) => {
                reject(err);
            });
        });
    }

    get(key) {
        return new Promise((resolve, reject) => {
            this.connect().then((client) => {
                client.get(key, (err, reply) => {
                    if (err) {
                        reject(err);
                    } else {
                        resolve(reply);
                    }
                });
            }).catch((err) => {
                reject(err);
            });
        });
    }
}

module.exports = RedisConfig;
```