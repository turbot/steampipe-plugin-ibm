# Table: ibm_is_flow_log

Flow Logs for VPC enable the collection, storage, and presentation of information about the Internet Protocol (IP) traffic going to and from network interfaces within your Virtual Private Cloud (VPC).

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log;
```

### List flow log collectors by name

```sql
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  name = 'steampipe01';
```

### List all inactive flow log collectors

```sql
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  not active;
```

### List all flow log collectors with auto delete disabled

```sql
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  not auto_delete;
```

### List flow logs with their corresponding VPC details

```sql
select 
  id, 
  name, 
  vpc ->> 'id' as vpc_id, 
  vpc ->> 'name' as vpc_name 
from 
  ibm_is_flow_log;
```