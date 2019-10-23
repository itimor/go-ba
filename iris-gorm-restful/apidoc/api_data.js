define({ "api": [
  {
    "type": "post",
    "url": "api/login",
    "title": "UserLogin",
    "name": "UserLogin",
    "group": "api",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>username of the User.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>password of the User.</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": " HTTP/1.1 200 OK\n {\n   \"status\": true,\n   \"msg\": \"sucess\",\n   \"data\": {\n      \"access_token\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzE4MTc5MTAsImlhdCI6MTU3MTgxNDMxMH0.5dAz2Fcfd1diaXzYONaehLB5tbf7Nyfa1HUGO3P4qew\"\n   }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": " HTTP/1.1 500 internal server error\n {\n   \"status\": false,\n   \"msg\": \"parames err\",\n   \"data\": null\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "controllers/auth.go",
    "groupTitle": "api"
  }
] });
