# R-0009: gson 探针战役 — 反射面实战能力边界报告

- 状态: 完成
- 关联: P-0011/P-0012（反射面）、P-0013（gson 战役）；ADR-0011
- 靶标: gson 2.8.9（Maven Central，258KB，214 类）
- 环境: darwin/arm64 · go1.26.5 · JDK 25 · 同机

## 探针结果矩阵

| 探针 | 内容 | Catty | JVM | 偏差 |
|---|---|---|---|---|
| P1 | forName from JAR + 构造 Gson | ✅ | ✅ | 无 |
| P2 | toJson(简单 POJO) | ✅ | ✅ | boolean 已修 |
| P3 | fromJson(json, Class) 反射反序列化 | ✅ | ✅ | 无 |
| P4 | 嵌套对象序列化 | ✅ | ✅ | 无 |
| P5 | @SerializedName 注解 | ✅ user_name/age_years | ✅ | **已修复** |
| P6 | 泛型 TypeToken<List<String>> | ❌ "Missing type parameter" | ✅ list_size=2 | 边界（需 ParameterizedType 模型） |

**结论：基本序列化/反序列化路径全通；@SerializedName 注解已支持；泛型 TypeToken 为唯一登记边界。**

## 新增内核表面清单（本轮 + 上轮合计）

### 标记接口 (10)
Cloneable, Closeable, Flushable, Serializable(已有), Type,
GenericArrayType, ParameterizedType, WildcardType, TypeVariable, Annotation

### 类 (20+)
java/lang/Class(正式定义+全套方法), java/lang/Enum(minimal),
java/lang/Number, java/lang/Float/Double(TYPE+parse),
java/lang/Void(TYPE), java/lang/ThreadLocal(single-slot),
java/util/concurrent/atomic.{AtomicInteger,AtomicBoolean,AtomicLong,
AtomicIntegerArray,AtomicLongArray}, java/util/BitSet,
java/util/LinkedHashMap, java/math/BigDecimal/BigInteger,
java/net/URL/URI/InetAddress, java/util/UUID/Currency/Locale/Calendar/
GregorianCalendar, java/lang/reflect/InvocationTargetException/
InstantiationException

### 方法 (40+)
Class: forName/getName/getSimpleName/isAssignableFrom/isInstance/
isArray/isEnum/isInterface/isPrimitive/getModifiers/getAnnotation(s)/
getSuperclass/getInterfaces/newInstance/cast/getField/getDeclaredField/
getDeclaredConstructor/getConstructor
Field: getName/getType/getGenericType/getDeclaringClass/getModifiers/
isSynthetic/isEnumConstant/setAccessible/isAccessible/get/set
Method: getName/invoke/getModifiers/equals/hashCode
Constructor: getName/newInstance/isAccessible/setAccessible/
equals/hashCode
String: format/clone
Collections: emptyList/emptySet/emptyMap/singletonList
System: getProperty(JDK8 合成)/format
ArrayList: addAll
LinkedHashMap: 全 HashMap 面
Object: getClass(已有)
数组伪类: clone()

## 过程中根修的内核 bug

| # | bug | 影响 |
|---|---|---|
| 1 | 验证器 ldc-class 推入目标类而非 java/lang/Class | gson Excluder isAssignableFrom 拒绝 |
| 2 | anewarray 组件描述符未规范化（缺 L…;） | 数组类身份分裂→CCE |
| 3 | 数组伪类无 methodsByKey | 数组 clone() NoSuchMethodError |
| 4 | Integer 缺 Number 父类 | gson 类型分派 CCE |
| 5 | hashKey 空 Fields/nil 键 panic | HashMap null key 崩溃 |
| 6 | primitiveClass 在 k.mu 下调 mustLookup → 死锁 | bootstrap 卡死 |
| 7 | verifier trailing locals 不兼容 | JsonReader 拒收 |
| 8 | popExpect 对 cat-2 top 值过严 | 同上 |
| 9 | stateCompatible stack 长度严格等号 | 同上 |

## 能力边界图（gson 2.8.9）

```
✅ 可用
├─ Gson() 构造
├─ toJson(POJO) 含嵌套对象
├─ fromJson(json, Class)
├─ JAR classpath 加载
├─ 反射字段枚举/get/set
└─ Method.invoke / Constructor.newInstance

⚠️ 受限（有 workaround 或降级）
├─ ~~@SerializedName → 已修复，输出与 JVM 一致~~
├─ getMethods 不含接口方法
└─ boolean 装箱走 Boolean.valueOf（已修）

❌ 不可用（精确边界）
├─ TypeToken<泛型> → 需要 ParameterizedType 完整模型
├─ 注解元数据保留/读取
├─ sun.misc.Unsafe 分配器路径
└─ EnumSet/EnumMap 等枚举容器
```

## 结论

反射最小面在 gson 2.8.9 的核心序列化/反序列化路径上实战可用。
泛型 TypeToken 是最大的单一缺口——需要实现 ParameterizedType 完整模型
（含 owner type、actual type arguments、raw type 三元组），估算 M 大小，
按真实需求排期。注解保留需要 classfile 解析层存储 RuntimeVisible/
InvisibleAnnotations 属性，估算 S-M。
