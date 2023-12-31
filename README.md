# Filter on Field
Allows for filtering Traefik Requests based on parameters that are passed in to the request matching a particular regular expression.

![Pipeline Status](https://gitlab.com/dominion-solutions-open-source/traefik-filter-on-field/badges/main/pipeline.svg)
![Latest Release]()

## Configuration Examples
### Docker
```yml
label:
- traefik.http.middlewares.<middlewareName>.traefikFilterOnField.fieldName='<fieldName>'
- traefik.http.middlewares.<middlewareName>.traefikFilterOnField.disallowedContent='<disallowedRegex1>,<disallowedRegex2>'
- traefik.http.middlewares.<middlewareName>.traefikFilterOnField.responseMessage="<Response Message>"
```

### File (YAML)
```yml
http:
  routers:
    <routerName>:
      service: <serviceName>
      middlewares:
      - <middlewareName>
  middlewares:
    <middlewareName>:
      plugin:
        traefikFilterOnField:
          fieldName: <fieldName>
          disallowedContent:
            - <disallowedRegex1>
            - <disallowedRegex2>
          responseMessage: <Response Message>
```

## Parameter Precedence
Parameters are approached in the following order:
1. **PUT** or **POST** Form Fields
2. **GET** Query String Fields.

:warning: If the same field exists in bothe a URL Query String and a Form Request, only the Form Request Field will be evaluated

## Disallowing Use Of A Parameter
```yml
middleware:
  filter-all-requests-with-disallowed-field:
    plugin:
      traefikFilterOnField:
        fieldName: disallowedField
        disallowedContent:
        # Regex that matches everything...
        - .*
```
