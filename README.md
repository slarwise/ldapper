# ldapsearch to json

Parse the output of ldapsearch, [LDAP Data Interchange Format](https://en.wikipedia.org/wiki/LDAP_Data_Interchange_Format), into json.

```sh
go build .

# Parse from stdin
ldapsearch -H ldap://hostname -LLL -b 'base_ou' 'memberOf=group1' displayName | ./ldapper

# Parse from file
./ldapper ldap_result.ldif
```

Example input:

```ldif
dn: cn=The Postmaster,dc=example,dc=com
objectClass: organizationalRole
cn: The Postmaster
```

Output:

```json
[
  {
    "cn": ["The Postmaster"],
    "dn": ["cn=The Postmaster,dc=example,dc=com"],
    "objectClass": ["organizationalRole"]
  }
]
```
