# marshalcache
MarshalCache序列化反序列化缓存
## 1.每个要缓存的对象提供方法：
`GetMcKey() string`
如果要序列化或者反序列化的对象没有这个实例方法，那么将无法应用缓存，而是直接调用对应的序列化反序列化方法。
## 2.缓存的写入
### 2.1Marshal的时候
mc.Marshal(obj)
只要满足GetMcKey接口的，那么获得Key，然后将序列化对象放入缓存
### 2.2Unmarshal的时候
如果Unmarshal的对象满足GetMcKey接口，那么就将传入的bytes放入缓存。
## 3子对象缓存（高级特性）

## 4性能测试结果
以JSON序列化为例，使用MarshalCache后性能明显提高
```
BenchmarkMarshalCacheImpl_Marshal_Direct
BenchmarkMarshalCacheImpl_Marshal_Direct-8     	    3181	    337121 ns/op
BenchmarkMarshalCacheImpl_Marshal_UseCache
BenchmarkMarshalCacheImpl_Marshal_UseCache-8   	 6956326	       169 ns/op
BenchmarkMarshalCacheImpl_Marshal_NoCache
BenchmarkMarshalCacheImpl_Marshal_NoCache-8    	    3555	    325300 ns/op
```