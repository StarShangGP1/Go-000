Week02 作业题目：
1. 我们在数据库操作的时候，比如 `dao` 层中当遇到一个 `sql.ErrNoRows` 的时候，是否应该 `Wrap` 这个 `error`，抛给上层。为什么？应该怎么做请写出代码？

个人理解：
1. 出现 `sql.ErrNoRows` 错误的时候，在 `dao` 层应该将错误信息的描述添加上去。
2. 在 `service` 层 进行 `wrap`，可避免重复定义 `error` 的报错处理。

