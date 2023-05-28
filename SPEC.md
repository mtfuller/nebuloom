# Specification

## `ROOT`
```
<COMPONENT>[]
```

## `COMPONENT`
```
<TYPE> <IDENTIFIER> {
    <COMPONENT_ELEMENT>[]
}
```

## `COMPONENT_ELEMENT`
```
<PROPERTY> | <COMPONENT>
```

## `PROPERTY`
```
<IDENTIFIER> -> <EXPR>
```

## `EXPR`
```
<VALUE> | <FUNC_CALL> | 
```

## `VALUE`
```
null | ".+" | [\d]+ | <OBJECT>
```

## `OBJECT`
```
{
    <PROPERTY>[]
}
```

## `FUNC_CALL`
```
<OBJECT_ACCESS> <OBJECT>?
```

## `OBJECT_ACCESS`
```
<IDENTIFIER>(.<OBJECT_ACCESS>)?
```