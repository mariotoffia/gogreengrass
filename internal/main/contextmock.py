from collections import namedtuple
import re
import warnings

ARN_FIELD_REGEX = \
    'arn:aws(?:-[a-z]+)*:lambda:([a-z]{2}(?:-[a-z]+)+-\d{1}):([0-9a-zA-Z]*):function:([a-zA-Z0-9-_]+)(?::(\$LATEST|[a-zA-Z0-9-_]+))?'


ClientContext = namedtuple(
    'ClientContext',
    ['client', 'custom', 'env']
)

ClientContextClient = namedtuple(
    'ClientContextClient',
    ['installation_id', 'app_title', 'app_version_name', 'app_version_code', 'app_version_package_name']
)

# ClientContextClient({
#                         'installation_id': 'apa',
#                         'app_title': 'testapp',
#                         'app_version_name': 'bubben',
#                         'app_version_code': 'älg',
#                         'app_version_package_name': 'buu'
#                         })
class _Context:
    def __init__(self, function_arn, invocation_id, client_context_encoded=None):
        arn_fields = FunctionArnFields(function_arn)

        self.invoked_function_arn = function_arn
        self.function_name = arn_fields.name
        self.function_version = arn_fields.qualifier
        self.aws_request_id = invocation_id

        if client_context_encoded:
            self.client_context = client_context_encoded
        else:
            self.client_context = ClientContext(
                    None,
                    {'subject': 'invoke/ggtest'},
                    None
                )

        self.identity = None

        # skip self.memory_limit_in_mb
        # skip self.log_group_name
        # skip self.log_stream_name

class FunctionArnFields:
    """
    @deprecated Since GGC version 1.9.0. Use the build_function_arn function instead.

    This class takes in a string representing a Lambda function's ARN (the qualifier is optional), parses that string
    into individual fields for region, account_id, name and qualifier. It also has a static method for creating a
    Function ARN string from those subfields.
    """
    @staticmethod
    def build_arn_string(region, account_id, name, qualifier):
        warnings.warn("Call to deprecated function build_arn_string.", DeprecationWarning, stacklevel=2)

        if qualifier:
            return 'arn:aws:lambda:{region}:{account_id}:function:{name}:{qualifier}'.format(
                region=region, account_id=account_id, name=name, qualifier=qualifier
            )
        else:
            return 'arn:aws:lambda:{region}:{account_id}:function:{name}'.format(
                region=region, account_id=account_id, name=name
            )

    @staticmethod
    def build_function_arn(unqualified_arn, qualifier):
        if qualifier:
            return '{unqualified_arn}:{qualifier}'.format(unqualified_arn=unqualified_arn, qualifier=qualifier)
        else:
            return unqualified_arn

    def __init__(self, function_arn_string):
        self.parse_function_arn(function_arn_string)

    def parse_function_arn(self, function_arn_string):
        regex_match = re.match(ARN_FIELD_REGEX, function_arn_string)
        if regex_match:
            region, account_id, name, qualifier = [s.replace(':', '') if s else s for s in regex_match.groups()]
        else:
            raise ValueError('Cannot parse given string as a function ARN.')

        self.region = region
        self.account_id = account_id
        self.name = name
        self.qualifier = qualifier
        self.unqualified_arn = function_arn_string[:-len(qualifier) - 1] if qualifier else function_arn_string

    def to_arn_string(self):
        return FunctionArnFields.build_function_arn(self.unqualified_arn, self.qualifier)
