# Table: ibm_certificate_manager_certificate

A Certificate Manager helps you to obtain, store and manage SSL/TLS certificates that you use for IBM Cloud deployments, or other cloud and on-premises deployments.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  issuer
from
  ibm_certificate_manager_certificate
where
  certificate_manager_instance_id = 'crn:v1:bluemix:public:cloudcerts:us-south:a/76aa4877fab6436db86f121f62faf221:f68bd88f-c4d4-4d81-9656-609e0b794c68::';
```

### List all imported certificates

```sql
select
  name,
  id,
  status,
  issuer,
  imported
from
  ibm_certificate_manager_certificate
where
  certificate_manager_instance_id = 'crn:v1:bluemix:public:cloudcerts:us-south:a/76aa4877fab6436db86f121f62faf221:f68bd88f-c4d4-4d81-9656-609e0b794c68::' and imported;
```
