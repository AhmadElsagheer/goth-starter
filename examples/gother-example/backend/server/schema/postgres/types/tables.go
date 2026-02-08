package sqtypes

import "github.com/bokwoon95/sq"

type ACCOUNTS struct {
	sq.TableStruct
	ID                    sq.UUIDField   `ddl:"type=uuid notnull primarykey default=gen_random_uuid()"`
	ACCOUNTID             sq.StringField `sq:"accountId" ddl:"type=text notnull"`
	PROVIDERID            sq.StringField `sq:"providerId" ddl:"type=text notnull"`
	USERID                sq.UUIDField   `sq:"userId" ddl:"type=uuid notnull references={users.id ondelete=cascade}"`
	ACCESSTOKEN           sq.StringField `sq:"accessToken" ddl:"type=text"`
	REFRESHTOKEN          sq.StringField `sq:"refreshToken" ddl:"type=text"`
	IDTOKEN               sq.StringField `sq:"idToken" ddl:"type=text"`
	ACCESSTOKENEXPIRESAT  sq.TimeField   `sq:"accessTokenExpiresAt" ddl:"type=timestamptz"`
	REFRESHTOKENEXPIRESAT sq.TimeField   `sq:"refreshTokenExpiresAt" ddl:"type=timestamptz"`
	SCOPE                 sq.StringField `ddl:"type=text"`
	PASSWORD              sq.StringField `ddl:"type=text"`
	CREATEDAT             sq.TimeField   `sq:"createdAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
	UPDATEDAT             sq.TimeField   `sq:"updatedAt" ddl:"type=timestamptz notnull"`
	_                     struct{}       `ddl:"index=\"userId\""`
}

type JWKS struct {
	sq.TableStruct
	ID         sq.UUIDField   `ddl:"type=uuid notnull primarykey default=gen_random_uuid()"`
	PUBLICKEY  sq.StringField `sq:"publicKey" ddl:"type=text notnull"`
	PRIVATEKEY sq.StringField `sq:"privateKey" ddl:"type=text notnull"`
	CREATEDAT  sq.TimeField   `sq:"createdAt" ddl:"type=timestamptz notnull"`
	EXPIRESAT  sq.TimeField   `sq:"expiresAt" ddl:"type=timestamptz"`
}

type SESSIONS struct {
	sq.TableStruct
	ID        sq.UUIDField   `ddl:"type=uuid notnull primarykey default=gen_random_uuid()"`
	EXPIRESAT sq.TimeField   `sq:"expiresAt" ddl:"type=timestamptz notnull"`
	TOKEN     sq.StringField `ddl:"type=text notnull unique"`
	CREATEDAT sq.TimeField   `sq:"createdAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
	UPDATEDAT sq.TimeField   `sq:"updatedAt" ddl:"type=timestamptz notnull"`
	IPADDRESS sq.StringField `sq:"ipAddress" ddl:"type=text"`
	USERAGENT sq.StringField `sq:"userAgent" ddl:"type=text"`
	USERID    sq.UUIDField   `sq:"userId" ddl:"type=uuid notnull references={users.id ondelete=cascade}"`
	_         struct{}       `ddl:"index=\"userId\""`
}

type USERS struct {
	sq.TableStruct
	ID                  sq.UUIDField    `ddl:"type=uuid notnull primarykey default=gen_random_uuid()"`
	NAME                sq.StringField  `ddl:"type=text notnull"`
	EMAIL               sq.StringField  `ddl:"type=text notnull unique"`
	EMAILVERIFIED       sq.BooleanField `sq:"emailVerified" ddl:"type=boolean notnull"`
	IMAGE               sq.StringField  `ddl:"type=text"`
	CREATEDAT           sq.TimeField    `sq:"createdAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
	UPDATEDAT           sq.TimeField    `sq:"updatedAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
	PHONENUMBER         sq.StringField  `sq:"phoneNumber" ddl:"type=text unique"`
	PHONENUMBERVERIFIED sq.BooleanField `sq:"phoneNumberVerified" ddl:"type=boolean"`
	ROLES               sq.JSONField    `ddl:"type=jsonb notnull"`
}

type VERIFICATIONS struct {
	sq.TableStruct
	ID         sq.UUIDField   `ddl:"type=uuid notnull primarykey default=gen_random_uuid()"`
	IDENTIFIER sq.StringField `ddl:"type=text notnull index"`
	VALUE      sq.StringField `ddl:"type=text notnull"`
	EXPIRESAT  sq.TimeField   `sq:"expiresAt" ddl:"type=timestamptz notnull"`
	CREATEDAT  sq.TimeField   `sq:"createdAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
	UPDATEDAT  sq.TimeField   `sq:"updatedAt" ddl:"type=timestamptz notnull default=CURRENT_TIMESTAMP"`
}
