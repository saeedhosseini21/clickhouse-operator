from clickhouse import * 
from kubectl import *
import settings 
from test_operator import *
from test_clickhouse import *  

from testflows.core import TestScenario, Name, When, Then, Given, And, main, run, Module, TE, args
from testflows.asserts import error

if main():
    with Module("main", flags=TE):
        with Given(f"Clean namespace {settings.test_namespace}"):
            kube_deletens(settings.test_namespace)
            kube_createns(settings.test_namespace)

        with Given(f"clickhouse-operator version {settings.version} is installed"):
            if kube_get_count("pod", ns='kube-system', label="-l app=clickhouse-operator") == 0:
                config = get_full_path('../deploy/operator/clickhouse-operator-install-template.yaml')
                kube_apply(f"<(cat {config} | "
                           f"OPERATOR_IMAGE=\"altinity/clickhouse-operator:{settings.version}\" "
                           f"OPERATOR_NAMESPACE=\"kube-system\" "
                           f"METRICS_EXPORTER_IMAGE=\"altinity/metrics-exporter:{settings.version}\" "
                           f"METRICS_EXPORTER_NAMESPACE=\"kube-system\" "
                           f"envsubst)", ns="kube-system")
            set_operator_version(settings.version)

        with Given(f"Install ClickHouse template {settings.clickhouse_template}"):
            kube_apply(get_full_path(settings.clickhouse_template), settings.test_namespace)

        with Given(f"ClickHouse version {settings.clickhouse_version}"):
            pass

        # python3 tests/test.py --only operator*
        with Module("operator", flags=TE):
            all_tests = [
                test_001,
                test_002,
                test_004,
                test_005,
                test_006,
                test_007,
                test_008,
                # (test_009, {"version_from": "0.8.0"}),
                (test_009, {"version_from": "0.9.6"}),
                test_010,
                test_011,
                test_011_1,
                test_012,
                test_013,
                test_014,
                test_015,
                test_016,
                test_017,
                test_018,
                test_019,
                test_020,
            ]
            run_tests = all_tests
            
            # placeholder for selective test running
            # run_tests = [test_020]

            for t in run_tests:
                if callable(t):
                    run(test=t, flags=TE)
                else:
                    run(test = t[0], args = t[1], flags=TE)

        # python3 tests/test.py --only clickhouse*
        with Module("clickhouse", flags=TE):
            all_tests = [
                test_ch_001,
            ]
        
            run_test = all_tests
            
            # placeholder for selective test running
            # run_test = [test_009]

            for t in run_test:
                run(test=t, flags=TE)
