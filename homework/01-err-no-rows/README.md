## 作业描述

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 回答

不必 wrap 这个错误抛给上层。

产生 `sql.ErrNoRows` 主要是出现在查询语句中，通常业务逻辑中，没有结果也算结果。

例如用户查看订单，没有更多订单需要返回给调用者的是*空数组*，而非错误。 

也可以认为 `sql.ErrNoRows` 是可处理的错误，根据「只处理一次错误」「之后不再报告当前错误」这两个原则不必再抛向上层。