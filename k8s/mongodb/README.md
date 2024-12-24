# MongoDB setup

## Setup for SPIFFE

```mongodb
use $external
db.runCommand({
  createUser: "CN=eventrunner,OU=threadr,O=carverauto,L=Carver,ST=MN,C=US",
  roles: [
    { role: "readWrite", db: "eventrunner" },
    { role: "userAdminAnyDatabase", db: "admin" }
  ]
})
```

```mongodb
use $external
db.getUsers()
```

```mongodb
$external> db.getUsers()
{
  users: [
    {
      _id: '$external.CN=eventrunner,OU=threadr,O=carverauto,L=Carver,ST=MN,C=US',
      userId: UUID('XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX'),
      user: 'CN=eventrunner,OU=threadr,O=carverauto,L=Carver,ST=MN,C=US',
      db: '$external',
      roles: [
        { role: 'readWrite', db: 'eventrunner' },
        { role: 'userAdminAnyDatabase', db: 'admin' }
      ],
      mechanisms: [ 'external' ]
    }
  ],
  ok: 1
}
```