{% set breadcumbs = [] %}
{% set breadcumbs_section = [] %}
{%- for node in schema.nodes_from_root -%}    
        {%- if not loop.first -%}
            {%- set _ = breadcumbs.append(node.name_for_breadcrumbs) -%}
            {%- if node.type_name == "object" -%}
                {%- set _ = breadcumbs_section.append(node.name_for_breadcrumbs) -%}
            {%- endif -%}
        {%- endif -%}
{%- endfor -%}

{% set section_name = breadcumbs_section|join('.')  %}
{% set variable_name = breadcumbs|last()  %}

**Example setting the default value** ({{schema.default_value}}):
```
{% if section_name != "" %}
[{{section_name}}]
{% endif %}
{{variable_name}}={{schema.default_value}}
```

