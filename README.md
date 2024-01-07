这是基于fabric-v2.2.0修改的支持国密算法的fabric，已通过命令行完成网络部署以及链码操作测试。

版权所有 上海旺链信息科技有限公司(http://www.vonechain.com)。


#### 一、修改说明


SM2为基于椭圆曲线密码的公钥密码算法标准，包含数字签名、密钥交换和公钥加密
SM3为密码哈希算法，用于替代MD5 / SHA-1 / SHA-256等国际算法。
SM4为分组密码，用于替代DES / AES等国际算法。
SM2，SM3，SM4现在是公开标准，直接使用[苏州同济区块链研究院实现的代码](https://github.com/tjfoc/gmsm)

通过国密SM2加解密、SM3withSM2签名验签，SM3哈希完成国密的改造.


1. BCCSP的全称是区块链密码服务提供者，它能够用来提供Fabric中加解密、签名校验相关功能 。
BCCSP通过Membership Service（成员服务提供者）给相关核心功能和客户端SDK提供加密算法相关的服务。
相关的核心功能集中在core中，包括consensus模块，endorser模块等。

[参考文章](https://blog.csdn.net/dyj5841619/article/details/90638054)

对国密算法的支持首先就要通过这个模块进行入手，创建bccsp的gm国密模块，实现思路是，参照sw模块来实现gm模块

(1)新增gmfactory.go，定义GMFactory结构体，增加gm文件夹，与sw类似

```markdown
bccsp/factory/gmfactory.go
bccsp/gm/certhelper.go
bccsp/gm/conf.go
bccsp/gm/dummyks.go
bccsp/gm/ecdsakey.go
bccsp/gm/fileks.go
bccsp/gm/hash.go
bccsp/gm/impl.go
bccsp/gm/internals.go
bccsp/gm/keyderiv.go
bccsp/gm/keygen.go
bccsp/gm/keyimport.go
bccsp/gm/sm2.go
bccsp/gm/sm2key.go
bccsp/gm/sm4.go
bccsp/gm/sm4key.go
```


(2)sw是bccsp中的软件实现方案，加入对sm2，sm4的处理，定义相关的结构体，修改清单如下

[参考文章](https://blog.csdn.net/dyj5841619/article/details/90639857)
```markdown
bccsp/sw/sm2.go //新增
bccsp/sw/sm2key.go //新增
bccsp/sw/sm4.go //新增
bccsp/sw/sm4key.go //新增
bccsp/sw/conf.go 
bccsp/sw/ecdsa.go
bccsp/sw/fileks.go
bccsp/sw/impl_test.go
bccsp/sw/keyderiv.go
bccsp/sw/keygen.go
bccsp/sw/keyimport.go
bccsp/sw/new.go
```

(3)factory是bccsp中的工厂服务，提供密码选项的配置，指定加解密方式，生成特定的服务工厂
```markdown
bccsp/factory/factory.go
bccsp/factory/factory_test.go
bccsp/factory/nopkcs11.go
bccsp/factory/nopkcs11_test.go
bccsp/factory/opts.go
bccsp/factory/opts_test.go
bccsp/factory/pkcs11.go

```
(4)utils是bccsp中的工具函数服务，提供一些支持服务，修改清单:
[参考文章](https://blog.csdn.net/dyj5841619/article/details/90640356)
```markdown
bccsp/utils/errs.go
bccsp/utils/io.go
bccsp/utils/keys.go
bccsp/utils/slice.go
```

(5)新增国密选项
```markdown
bccsp/hashopts.go
bccsp/opts.go
```
(6)修改signer包中证书的解析函数
```markdown
bccsp/signer/signer.go
```

2. idemix（Identity Mixer）的核心是零知识证明（Zero Knowledge Proof）。
用户无需暴露私有数据以及任何有用的信息，也能证明自己拥有这些私有数据，对方能够进行有效验证，这就是零知识证明。
idemix是一个密码协议套件（X.509+加密算法），保留隐私实现匿名性，交易时不用透露交易者的身份，而且交易间是无关联的，不可往前追溯。


[参考文章](https://blog.csdn.net/dyj5841619/article/details/90638997)


idemix包中的修改需要把所有ecdsa算法改为sm2算法，修改的文件清单如下：
```markdown
bccsp/idemix/bridge/bridge_test.go
bccsp/idemix/bridge/revocation.go
bccsp/idemix/bridge/signaturescheme.go
bccsp/idemix/handlers/idemix.go
bccsp/idemix/handlers/mock/revocation.go
bccsp/idemix/handlers/mock/signature_scheme.go
bccsp/idemix/handlers/revocation.go
bccsp/idemix/handlers/revocation_test.go
```


3. common包是一些通用的模块的集合，主要是包含各种策略定义工具包、编译工具包、配置文件工具包、加解密工具包、账本包和docker挂载相关的工具包、基本数据metadata工具包等


主要修改crypto包里面的ca.go和key.go，tools包作相应调整和适配


4. 将crypto/x509替换成github.com/tjfoc/gmsm/sm2，
将ecdsa替换成github.com/tjfoc/gmsm/sm2，将crypto/tls替换成github.com/tjfoc/gmtls，将google.golang.org/grpc/credentials替换成github.com/tjfoc/gmtls/gmcredentials
修改的列表：

(1)core包：大部分核心实现代码都在本包下，其他包的代码封装上层接口，最终调用本包内代码
```markdown
core/chaincode/accesscontrol/mapper.go
core/chaincode/lifecycle/serializer.go
core/deliverservice/config.go
core/deliverservice/deliveryclient.go
core/endorser/msgvalidation.go
core/operations/system.go
core/operations/tls.go
core/peer/config.go
```

(2)gossip包：实现gossip协议

```markdown
gossip/api/crypto.go
gossip/comm/comm_impl.go
gossip/comm/crypto.go
gossip/comm/crypto_test.go
gossip/util/grpc.go
```

```markdown
idemix/revocation_authority.go
idemix/signature.go
```
(3)msp包：Member Service Provider包

```markdown
msp/cert.go
msp/cert_test.go
msp/configbuilder.go
msp/identities.go
msp/mgmt/mgmt.go //主要文件，localMsp和mspMap都在这个文件，还有多个管理函数
msp/mgmt/testtools/config_test.go
msp/mspimpl.go //实现MSP接口，结构为bccspmsp
msp/mspimplsetup.go
msp/mspimplsetup_test.go
msp/mspimplvalidate.go
```
(4)order包：order服务相关的入口和框架代码

```markdown
orderer/common/cluster/comm.go
orderer/common/cluster/comm_test.go
orderer/common/cluster/connections.go
orderer/common/cluster/deliver.go
orderer/common/cluster/rpc_test.go
orderer/common/cluster/util.go
orderer/common/cluster/util_test.go
orderer/consensus/etcdraft/membership.go
orderer/consensus/etcdraft/util.go
orderer/consensus/kafka/config.go
```

5. sha256替换成sm3

```markdown
protoutil/txutils.go
```


6. 重点说明

(1)证书相关的结构体的修改，增加sm2证书和sm2私钥
```markdown
type CA struct {
	Name               string
	Country            string
	Province           string
	Locality           string
	OrganizationalUnit string
	StreetAddress      string
	PostalCode         string
	Signer             crypto.Signer
	SignCert           *x509.Certificate
	SignSm2Cert        *sm2.Certificate //sm2证书
	Sm2Key             *sm2.PrivateKey //sm2私钥
}

```
(2)证书生成的函数修改，指定证书的签名方式和计算私钥的SKI
```markdown

func NewCA(
	baseDir,
	org,
	name,
	country,
	province,
	locality,
	orgUnit,
	streetAddress,
	postalCode string,
) (*CA, error) {
	var ca *CA
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, err
	}
	priv, err := csp.GeneratePrivateKey(baseDir)
	if err != nil {
		return nil, err
	}
	template := x509Template()
	//this is a CA
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageDigitalSignature |
		x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign
	template.ExtKeyUsage = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageServerAuth,
	}
	//set the organization for the subject
	subject := subjectTemplateAdditional(country, province, locality, orgUnit, streetAddress, postalCode)
	subject.Organization = []string{org}
	subject.CommonName = name
	template.Subject = subject
	templateSm2 := gm.ParseX509CertificateToSm2(&template)
	//TODO important
	templateSm2.SubjectKeyId = computeSKI(priv)  //指定证书的SKI
	sm2PubKey := priv.PublicKey
	if err != nil {
		errors.Errorf("error,%v", err)
	}
	templateSm2.SignatureAlgorithm = sm2.SM2WithSM3  //指定证书的签名算法
	sm2Cert, err := genCertificateSM2(
		baseDir,
		name,
		templateSm2,
		templateSm2,
		&sm2PubKey,
		priv,
	)
	if err != nil {
		return nil, err
	}
	ca = &CA{
		Name: name,
		Signer:             priv,
		Country:            country,
		Province:           province,
		Locality:           locality,
		OrganizationalUnit: orgUnit,
		StreetAddress:      streetAddress,
		PostalCode:         postalCode,
		SignSm2Cert:        sm2Cert,
		Sm2Key:             priv,
	}
	return ca, err
}

//定义sm2私钥的ski计算函数
func computeSKI(privKey *sm2.PrivateKey) []byte {
	// Marshall the public key
	raw := elliptic.Marshal(privKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)
	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

```

(3)sm2私钥生成
```markdown
//TODO  SM2
func GeneratePrivateKey(keystorePath string) (*sm2.PrivateKey, error) {

	priv, err := sm2.GenerateKey()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate private key")
	}
	pkcs8Encoded, err := sm2.MarshalSm2PrivateKey(priv,nil) //解析sm2证书方式
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal private key")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8Encoded})
	keyFile := filepath.Join(keystorePath, "priv_sk")
	err = ioutil.WriteFile(keyFile, pemEncoded, 0600)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to save private key to file %s", keyFile)
	}
	return priv, err
}


```

#### 二、编译

进入你克隆到本地的fabric仓库的根目录

```markdown
go mod vendor
```
```sql
cd vendor/github.com/tjfoc/gmtls   //暂时如此操作,防止后面的问题
go mod vendor
```

```sql
make docker //回到克隆到本地的fabric仓库的根目录
```

#### 三、问题

1.出现的问题
```markdown
remote error: tls: unexpected message
```

解决办法
```sql
cd vendor/github.com/tjfoc/gmtls
go mod vendor
```

2.出现的问题：找不到rootca的证书链

解决办法
```sql
templateSm2.SubjectKeyId = computeSKI(priv)  //指定证书的SKI
```

3.出现的问题：验证签名失败

解决办法
```sql
templateSm2.SignatureAlgorithm = sm2.SM2WithSM3  //指定证书的签名算法
```


#### 四、启动
```markdown
docker-compose -f docker-compose-cli.yaml -f docker-compose-etcdraft2.yaml -f docker-compose-ca.yaml  up -d
```



