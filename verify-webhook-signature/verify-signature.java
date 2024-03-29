import java.security.KeyFactory;
import java.security.PublicKey;
import java.security.Signature;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

public class Application {

    public static void main(String[] args) throws Exception {

        String encodedSandboxPublicKey = "-----BEGIN PUBLIC KEY-----" +
                "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwpb91cEYuyJNQepZAVfP" +
                "ZIlPZfNUefH+n6w9SW3fykqKu938cR7WadQv87oF2VuT+fDt7kqeRziTmPSUhqPU" +
                "ys/V2Q1rlfJuXbE+Gga37t7zwd0egQ+KyOEHQOpcTwKmtZ81ieGHynAQzsn1We3j" +
                "wt760MsCPJ7GMT141ByQM+yW1Bx+4SG3IGjXWyqOWrcXsxAvIXkpUD/jK/L958Cg" +
                "nZEgz0BSEh0QxYLITnW1lLokSx/dTianWPFEhMC9BgijempgNXHNfcVirg1lPSyg" +
                "z7KqoKUN0oHqWLr2U1A+7kqrl6O2nx3CKs1bj1hToT1+p4kcMoHXA7kA+VBLUpEs" +
                "VwIDAQAB" +
                "-----END PUBLIC KEY-----";

        String webhookBody = "{\"data\":{\"resource\":{\"id\":49983981,\"profile_id\":16055450,\"account_id\":14124090,\"type\":\"transfer\"},\"current_state\":\"incoming_payment_waiting\",\"previous_state\":null,\"occurred_at\":\"2021-08-23T10:12:50Z\"},\"subscription_id\":\"90aa8e14-4ef1-4a56-861c-f3c9cde097ea\",\"event_type\":\"transfers#state-change\",\"schema_version\":\"2.0.0\",\"sent_at\":\"2021-08-23T10:12:50Z\"}";
        String signatureHeader = "wKcKCYXAzxNgiu7xmoDm943NUni7Rz33QN8JkEA9dWSGebgndonabgSj18Y4C08OrwVmueGsED2s00M7DtJVcYKOS1i3G4TMVx+mgM3aL9djMBkQtiYNBFUd6wrPI7ZUNHv/TrlKSjTMc+6JFvUvJ7owY3z85e3I4jLRLJowMFvO8kvCJ60+1pY9wDwZvtZ//WS93LrwGjk9Dvwzpmu0w+P4J75tETT5qC3Uv0y5G2yO8SEoO3yNP/tg/BOli02niHb53vEOUWUb9bly6thnfMoXoiV/osoGxgF20R58RlvkAmezyyl1Sv542TfS2DpiwVnmjjjkCyXeSUcKookYLQ==";

        RSAPublicKey publicKey = parsePKCS8PublicKey(encodedSandboxPublicKey);
        Signature sign = Signature.getInstance("SHA256WithRSA");
        sign.initVerify(publicKey);
        sign.update(webhookBody.getBytes());

        byte[] signatureBytes = Base64.getDecoder().decode(signatureHeader);

        var isVerified = sign.verify(signatureBytes);
        System.out.println("Is signature verified? " + isVerified);
    }

    public static RSAPublicKey parsePKCS8PublicKey(String publicKey) throws NoSuchAlgorithmException, InvalidKeySpecException {
        var publicKeyWithoutHeaders = publicKey
                .replace("-----BEGIN PUBLIC KEY-----", "")
                .replace("-----END PUBLIC KEY-----", "")
                .replace("\n", "")
                .trim();

        X509EncodedKeySpec keySpec = new X509EncodedKeySpec(Base64.getDecoder().decode(publicKeyWithoutHeaders));
        KeyFactory keyFactory = KeyFactory.getInstance("RSA");
        return (RSAPublicKey) keyFactory.generatePublic(keySpec);
    }
}