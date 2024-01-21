
# Rest API Raw model

## 1. Filtering / Find

### A. Comparison Operator

filter have 3 field alias: `filter`, `where`, and `find`

- EQ : This operator of "=" or same
```
{
  "find": {
    "op": "$eq",
    "field": "name",
    "value": "Todo"
  }
}
```

<!-- - NOT: This operator of "<>" or "Not" or not same
```
"find":{
  "op": "$not",
  "field": "name",
  "value": "Todo"
}
``` -->

- LIKE: This operator of "LIKE" or regular expression for search
```
{
  "find"{
    "op": "$like",
    "field": "name",
    "value": "Todo"
  } 
}
```

- ILIKE: This operator of "ILIKE" or irregular expression for search
```
{
  "find"{
    "op": "$like",
    "field": "name",
    "value": "Todo"
  } 
}
```

- GT: This operator of ">" or Greater
```
{
  "find":{
    "op": "$gt",
    "field": "name",
    "value": "Todo"
  }
} 
```

- GTE: This operator of ">=" or greater than
```
{
  "find":{
    "op": "$gte",
    "field": "name",
    "value": "Todo"
  } 
}
```

- LT: This operator of "<" or Lower
```
{
  "find":{
    "op": "$lt",
    "field": "name",
    "value": "Todo"
  } 
}
```

- LTE: This operator of "<=" or Lower Than
```
{
  "find":{
    "op": "$lte",
    "field": "name",
    "value": "Todo"
  } 
}
```

- NIN: This operator of "Not IN"
```
{
  "find":{
    "op": "$nin",
    "field": "name",
    "value": ["Todo", "Foo"]
  } 
}
```

- IN: This operator of "IN"
```
{
  "find":{
    "op": "$in",
    "field": "name",
    "value": ["Todo", "Foo"]
  }
}
```

### B. Conjunction Operators

- AND
```
{
  "find":{
    "op": "$and",
    "items": [
      {
        "op": "eq",
        "field": "name",
        "value": "Todo"
      },
      {
        "op": "eq",
        "field": "name",
        "value": "Foo"
      } 
    ]
  }
}
```

- OR
```
{
  "find":{
    "op": "or",
    "items": [
      {
        "op": "$eq",
        "field": "name",
        "value": "Todo"
      },
      {
        "op": "$eq",
        "field": "name",
        "value": "Foo"
      }
    ]
  }
}
```

- NOT
```
{
  "find":{
    "op": "$not",
    "items": [
      {
        "op": "$eq",
        "field": "name",
        "value": "Todo"
      },
      {
        "op": "$eq",
        "field": "name",
        "value": "Foo"
      }
    ]
  }
}
```

- NOR
```
{
  "find":{
    "op": "$nor",
    "items": [
      {
        "op": "$eq",
        "field": "name",
        "value": "Todo"
      },
      {
        "op": "$eq",
        "field": "name",
        "value": "Foo"
      }
    ]
  }
}
```

```
NOTE:
op = Operator
field = Field in model
items = for item conjuction operator
value = value for field
```

## 2. Sorting
sort have 3 field alias : `sort`, `sortby`, and `orderby`

sorting fields list separated by comma (",")

Example:
```
{
  "sortby":["id asc", "name asc"]
}
```

sort have function `field()` for sort manually with custom value
Example:
```
{
  "sortby":["field(userid, 1, 2) asc", "name asc"]
}
```
API field : field(fieldname, ...value)

## 3. Limit

Example:
```
{
  "limit":10
}
```

## 4. Skip
sort have 2 field alias : `skip`, and `offset`
```
{
  "offset":10
}
```

## 5. Select
select for specific field output if from sql use `map[string]interface{}`
```
{
  "selects":["name", "email"]
}
```

