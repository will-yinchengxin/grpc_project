报错: 2023/11/15 18:50:58 rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: ALTS: untrusted platform. ALTS is only supported on GCP"

> ALTS认证只在Google Cloud Platform (GCP) 上受支持，无法在其他平台上使用。ALTS是一种特定于GCP的认证协议，用于在GCP环境中进行gRPC通信的安全认证。
> 
> 如果在本地环境或其他云平台上运行该示例，将无法使用ALTS认证。可以考虑使用其他认证方式，如TLS认证或基于令牌的认证，以满足需求。

