# GoLeap
It's a simple MySQL compatible ORM for Go, it makes it easy to manage your database with structs.

## TODO

```SQL
 SELECT `t0`.`id` FROM `test`.`simple` `t0` LIMIT 1;
```

```bash
2023/02/05 15:28:00 Extra.BaseModel.Recursive.Extra.ExtraJump
2023/02/05 15:28:00 add to visited: test@base::extra.extra_id/test@extra::extrajump.extra_jump_id

2023/02/05 15:28:00 Extra.BaseModel.Extra.ExtraJump
2023/02/05 15:28:00 field already visited: test@base::extra.extra_id/test@extra::extrajump.extra_jump_id
```

## SOLUTION

~~Try to use visited fields in scoped schema. Not sure if it's the best solution, but it can work.~~

~~Try with original type, maybe.~~

Recursive.Extra.BaseModel.Recursive.Extra.ExtraJump.JumpToBase

Try to count unique schema names, and use it when generate key.

