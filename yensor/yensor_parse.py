#!/usr/bin/env python

import yaml

def parse_yaml_from_path(path: str) -> dict:
    # return python dict from yaml path
    with open(path, 'r') as stream:
        try:
            y = yaml.load(stream)
            return y
        except yaml.YAMLError as exc:
            print(exc)
            return None


def create_model_and_arch_config(path: str) -> (dict, dict):
    # return the model and archtitecture configuration dicts
    model_config = parse_yaml_from_path(path)
    # create architecture config 
    if model_config['architecture']['yaml']:
        arch_config = parse_yaml_from_path(model_config['architecture']['yaml'])
    else:
        arch_config = model_config['architecture']
    
    return (model_config, arch_config)

MC, AC = create_model_and_arch_config("./model_config.yaml")
print(MC)
print(AC)