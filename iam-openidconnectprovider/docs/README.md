# Smaato::IAM::OpenIDConnectProvider

Create OpenID Connect (OIDC) identity provider.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Smaato::IAM::OpenIDConnectProvider",
    "Properties" : {
        "<a href="#url" title="Url">Url</a>" : <i>String</i>,
        "<a href="#clientidlist" title="ClientIDList">ClientIDList</a>" : <i>[ String, ... ]</i>,
        "<a href="#thumbprintlist" title="ThumbprintList">ThumbprintList</a>" : <i>[ String, ... ]</i>,
    }
}
</pre>

### YAML

<pre>
Type: Smaato::IAM::OpenIDConnectProvider
Properties:
    <a href="#url" title="Url">Url</a>: <i>String</i>
    <a href="#clientidlist" title="ClientIDList">ClientIDList</a>: <i>
      - String</i>
    <a href="#thumbprintlist" title="ThumbprintList">ThumbprintList</a>: <i>
      - String</i>
</pre>

## Properties

#### Url

The URL of the OpenID Connect (OIDC) identity provider. The URL must begin with https:// and should correspond to the iss claim in the provider's OpenID Connect ID tokens.

_Required_: Yes

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>255</code>

_Pattern_: <code>^(https)://.*$</code>

_Update requires_: [Replacement](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-replacement)

#### ClientIDList

A list of client IDs (also known as audiences).

_Required_: Yes

_Type_: List of String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ThumbprintList

A list of server certificate thumbprints for the OpenID Connect (OIDC) identity provider's server certificates.

_Required_: No

_Type_: List of String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the Name.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Name

The identifier of new IAM OpenID Connect provider that is created. The name is the OpenID Connect URL without https://.

#### Arn

The Amazon Resource Name (ARN) of the new IAM OpenID Connect provider that is created.

