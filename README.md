# GREST - Next Generation API Framework

The main goal of GREST is to make possible the implementation of an RESTFul APIs that consumes an SQL database without any programming skills.

In present days, an API that follows the REST principles is made defining some constraints and implementing a way to parse request contents in a language that SQL databases understand.

![REST API Overview](/img/img1.png)

GREST was already built with this implementations and constraints, it's up to the business administrator only the task of customize the prebuilt functionalities to fit it's integration with the application.

---

## GREST benefits:

### 1. Eficiency

By utilize Golang, a compiled language, and by make use of good algorithms of processing, GREST is capable to perform an average of only 10ms of processing time.

### 2. Scalabilty

Due to minimal acoplation, it's possible to configure pools of processing, make use of load balancers, automate his functioning (yes, the result API has his own API), as also integrate it with other systems.

### 3. Security

Included in the constrains there is a dozens of security parameters that can be customized, including WAF (Web Application Firewall), rate limits, authentication, authorization, etc.

### 4. Productivity

In last, but definetly not less important, there is the time necessary to get the things working, GREST use makes hours of development and integration unecessary. Some configurations and you are ready to **Go** :)

---

## Getting started

### First things first

The first step is to understand *how* GREST works, see the image below:

![GREST Pivot](/img/img2.png)

The pivot of GREST functioning is the concept of `behavior`. The behavior is how the data will be parsed in a communication with the database. 

That behaviors (yes, it's possible to have many) are composed of two things: `path mapping` and `key mappings`, that are only key value maps.

Behaviors are composed of only one path mapping and of a list of key mappings.

#### Path mappings

Path mappings tells to GREST how to provide the access of an specific group of data to the world. In other worlds, it's responsability is to tell how people access your database tables.

##### Example

`"/brazil/library/books"` -------> `"books_of_brazil_library"`

The path mapping above says: Hey GREST! Clients will access my `table` in this `endpoint`.

#### Key Mappings

Key mappings tells to GREST how and which parts of this group of data will be accessed. In other words, it's responsabilty is to tell how people will reference the columns of the tables that they are acessing.

##### Example

`"id"` -------> `"isbn_number"`

The key mapping above says: Hey GREST! Clients need to pass this `name` when they want to make things with the values in this `column`.

---

### Usage

Having you in hands an active database, maybe containing something like this:

*Employees*

| Employee_name | State | Telephone | E-mail |
| ---- | ---- | --------- | ------ |
| Joseph | SC | (99) 91111-1111 | jo@email.com |
| Pietro | SP | (99) 92222-2222 | pi@email.com |
| Matheus | SP | (99) 9333-3333 | ma@email.com |
| Paul | SP | (99) 94444-4444 | pa@email.com |

To make it visible to world, in `/business/our-people` endpoint, do the following:

1. To define the path mapping

```
curl -s http://[IP of server running GREST]:9090/config/path-mappings -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "id": "100"
    "path": "/business/our-people",
    "table": "Employees"
  }
}

```

2. To define the key mappings


```
curl -s http://[IP of server running GREST]:9090/config/key-mappings -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "id": "201"
    "key": "name",
    "column": "Employee_name"
  }
}

```

```
curl -s http://[IP of server running GREST]:9090/config/key-mappings -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "id": "202"
    "key": "email",
    "column": "E-mail"
  }
}

```

3. To group it in a behavior

```
curl -s http://[IP of server running GREST]:9090/config/behaviors -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "path_mapping_id": "101",
    "key_mapping_id": "201"
  }
}

```

```
curl -s http://[IP of server running GREST]:9090/config/behaviors -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "path_mapping_id": "101",
    "key_mapping_id": "202"
  }
}

```

Data defined in this moment in GREST:

```
{
  "path-mapping": {
    "path": "/business/our-people",
    "table": "Employees"
  },
  "key-mappings": [
    {
      "key": "name",
      "column": "Employee_name"
    },
    {
      "key": "email",
      "column": "E-mail"
    }
  ]
}
```

---

At this point, you already will be able to access the `Employee_name` and the `E-mail` column from `Employees` table stored in your database:


```
curl -s http://[IP of server running GREST]:8080/business/our-people
```

You also will be able to make request with some filters, like:

```
curl -s http://[IP of server running GREST]:8080/business/our-people?name=Paul
```


The response payload will always follow the following format:

```
{
  "status": [ status code ],
  "response": [ data retrieved from database ],
  "errors": [ erros that occurred in the process ]
}

```

---

#### Requests

GREST supports 6 different HTTP methods: GET, POST, PUT, DELETE, HEAD, OPTIONS. 

Im below, examples of all:


**GET**

```
curl -s http://[IP of server running GREST]:8080/business/our-people?name=Paul&email=pa@email.com
```

**POST**

```
curl -s http://[IP of server running GREST]:8080/business/our-people -X POST -d [payload]
```

```
Payload:

{
  "Set": {
    "name": "John",
    "email": "john@gmail.com"
  }
}
```

**PUT**

```
curl -s http://[IP of server running GREST]:8080/business/our-people -X PUT -d [payload]
```

```
Payload:

{
  "Must": {
    "name": "Paul"
  }
  
  "Set": {
    "email": "paul@gmail.com"
  }
}
```

**DELETE**

```
curl -s http://[IP of server running GREST]:8080/business/our-people -X DELETE -d [payload]
```

```
Payload:

{

  "Must": {
    "name": "Paul"
  }
  
}
```

**HEAD**

```
curl -s http://[IP of server running GREST]:8080/business/our-people
```

> HEAD method request and response does not contains body


**OPTIONS**

```
curl -s http://[IP of server running GREST]:8080/business/our-people -X OPTIONS
```

> OPTIONS method request and response does not contains body




 





