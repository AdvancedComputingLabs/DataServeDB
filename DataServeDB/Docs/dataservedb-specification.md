# DataserveDB Software Design Specification

1. ## Authentication and User Privileges
    > **NOTE**: Use hash(salt + password); salt in the beginning makes it harder to crack.

    1.1 **Authentication Schemes:**
    
      1.1.1 **CLSimpleAuth1**
      * Uses sha3-256 ( salt + password ) to store password.
      * Only salt and hash( salt + password ) is stored on the server side.
      * User sends passwords and server converts to hash to compare on the server side.
      * **Pros:**
        * Passwords are not stored on the server side in plain text.
        * Compares are fast and can serve lot of requests.
      * **Cons:**
        * No expiry on the hash compare; although server can be set to force user to change password after specified period.
        * Password is sent in plain text; but it is over secure connection like ssl then it is safe against evasdropping.
      
    1.2. **User Privileges Format and Access Codes:**
      * User authentication details are stored in hash table with user name and 'userauthobject' or hastable of pointers to 'userauthobjects'.
      * User auth object contains user's hashed password and claims.
      * Structure:
        ```go
        //global UserAuthObject1 wil be nil if user is in database(s).
        DbsUsers hastable[username:discrimnated union{global UserAuthObject1 or database level access hashtable[dbname:UserAuthObject1]}]
        
        UserAuthObject1 {
            AuthScheme string
            PwdH string //Password Hash in the format for the auth scheme.
            PwdRenewalIntervelWks uint //Interval is in weeks, 1 means after every 1 week; 0 means disabled.
            PwdLastChanged db datetime //Zero means not set and if there is renewal interval then user will be asked change password at first login.
            IsDbsRootUser bool //Dbs means database server.
            Claims hashtable {
                Claims hashtable { ... }
            }
        }
        ```
      * Claims can contain claims, it is hierarchical structure.
      * Standard claims are as follows:
        * Owner: has full access to the object: db, table, or any other db object.
        * Reader
        * Writer
        * Names are qualified with dot notation. Example: hastable[dbname.tablename:claims].
      * Database server default root:
        * Default root user is created at installation time.
        * Root user has full access to the database system and all the databases; it is meant for database server systems administrators.
        * Database users are stored in database, so if database is moved then database users are moved with it.
            * Conflict resolution with main dbs hashtable:
                * Main hashtable can contain pointers to multiple sub user auth objects.
                * When there are users with same name in mutiple databases, main hashtable should combine with userauthobject pointers hashtable.
                * Database level username cannot have global userauthobject; it is to avoid conflicts. If there exist a same username with global object it should error with conflict when attaching the database.
                * Implementation and optimization details are left to the implementator.
        * Authentication requires object to know if authentication data was sent over secure connection and type of secure connection. This way some operations that needs secure connection can be forced to be done on secure connection and they can limit the type of secure like minimum version of ssl to add furture security to the communication.