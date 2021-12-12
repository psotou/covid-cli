# COVID CLI

Covid CLI retorna la data de los casos COVID a nivel nacional, regional y comunal. Para los casos nacionales y regionales, retorna el número de casos por día; para los casos comunales, retorna el número de casos acumulados a la fecha para aquellas fecha que poseen registros.

Todos los datos son extraídos del repositorio del Ministerio de Ciencias. Para los casos nacionales uso el [producto 5](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto5), para los regionales el [producto 4](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto4) y para los comunales, el [producto 1](https://github.com/MinCiencia/Datos-COVID19/tree/master/output/producto1)

## Ejemplo de uso

Asumiendo que se compila este repo en local y que se nombra el archivo ejecutable **`covid`**:

```bash
covid <command> \
    --days <número de los días a consultar > \
    --name <nombre de región o comuna según el command>
```

Donde `command` &in; {`region`, `comuna`}. 

El CLI tiene como defecto 7 días, región metropolitana y comuna ñuñoa. Por lo que si se ejecuta:

```bash
covid region
```

Devolverá los datos para los últimos 7 días de la región metropolitana. Y si se ejecuta:

```bash
covid comuna
```

Devolverá las últimas fechas que tengan datos dentro de los últimos 7 días para la comuna de ñuñoa.

Ahora bien, si se quiere conocer los datos comunales para la comuna de Ñiquén (regióndel Ñuble) durante los últimos 14 días, se debe hacer:

```bash
covid comuna --days 12 --name niquen
```

Y para conocer lo datos de los últimos 12 días para la región del Ñuble:

```bash
covid región --days 12 --name nuble
```

**Nota:** no importa el orden de lo subcomando, es decir, bien podría ir primero `--name`. Además, se puede usar con un guión (por ejemplo, `-name`) si se prefiere.