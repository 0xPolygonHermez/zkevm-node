{% set my_string = [] %}
{%- for node in schema.nodes_from_root -%}    
        {%- if not loop.first -%}
            {%- set _ = my_string.append(node.name_for_breadcrumbs) -%}
        {%- endif -%}
{%- endfor -%}


{%- filter md_escape_for_table -%}
{#
  {%- if config.show_breadcrumbs -%}
    {%- for node in schema.nodes_from_root -%}
      {{ node.name_for_breadcrumbs }}{%- if not loop.last %} > {% endif -%}
    {%- endfor -%}
  {%- else -%}
    {{ schema.name_for_breadcrumbs }}
  {%- endif -%}
#}
{%- if schema.type_name == "object" -%}
[{{ my_string|join('.') }}]
{%- else -%}
{{ my_string|join('.') }}
{%- endif -%}
{%- endfilter -%}
