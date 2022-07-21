# COVID CLI

Covid CLI retorna la data de los casos COVID a nivel nacional, regional y comunal. Para los casos nacionales y regionales, retorna el número de casos por día; para los casos comunales, retorna el número de casos acumulados para aquellas fechas que poseen registros (este dato se actualiza cada 3 ó 4 días).

Todos los datos son extraídos del repositorio del Ministerio de Ciencias. Para los casos nacionales uso el [producto 5](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto5), para los regionales el [producto 4](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto4) y para los comunales, el [producto 1](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto1)

## Uso

Asumiendo que se clona este repo y que copia el ejecutal de la ruta `bin/covid` en, digamos, `/usr/local/bin/`, el uso del CLI sería:

```bash
covid <option> \
    -days <número de los días a consultar > \
    -name <nombre de región o comuna según el option>
```

Donde `option` &in; {`region`, `comuna`}. 

El CLI tiene como defaults 7 días, región Metropolitana y comuna Ñuñoa. Por lo que si se ejecuta:

```bash
covid region
```

Devolverá los datos para los últimos 7 días de la región metropolitana. Y si se ejecuta:

```bash
covid comuna
```

Devolverá las últimas fechas que tengan datos dentro de los últimos 7 días para la comuna de Ñuñoa.

## Convenciones para `-name` 

### Regiones

Los nombres de las regiones están mapeadas al nombre con los siguientes valores:

```bash
arica
tarapaca
antofagaste
atacama
coquimbo
valparaiso
metropolitana
ohiggins
maule
nuble
biobio
araucania
losrios
loslagos
aysen
magallanes
```

De esta manera, para obtener los datos regionales de los últimos 12 días de la región de Los Ríos corremos:

```bash
covid region -days 12 -name losrios
```

### Comunas

Los nombres de las comunas se deben escribir con minúscula y sin caracteres especiales (tildes o eñes). Para los nombres de comunas separados por espacios, el nombre se debe escribir como string. Como ejemplo del primer caso, tomemos la comuna de Ñiquén, en donde corremos:

```bash
covid comuna -days 12 -name niquen
```

Para el segundo caso, tomemos la comuna de San Carlos, en donde corremos:

```bash
covid comuna -days 12 -name 'san carlos'
```

**Nota:** no importa el orden del subcomando, es decir, bien podría ir primero `-name`.

## Usage

Regiones:

![](/static/region.gif)

Comunas:

![](/static/comuna.gif)

## TO-DO

+ [ ] Mejorar errores
+ [ ] Agregar tests!
+ [ ] Agegar Github Action
+ [ ] Agregar paquetes para win y linux
<!-- dummy commit: on -->
