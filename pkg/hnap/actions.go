package hnap

const (
	// Confirmed actions

	GetMultipleHNAPs = "GetMultipleHNAPs"

	GetHomeAddress = "GetHomeAddress"
	GetHomeConnection = "GetHomeConnection"


	GetMotoLagStatus = "GetMotoLagStatus"
	GetMotoStatusConnectionInfo = "GetMotoStatusConnectionInfo"
	GetMotoStatusDownstreamChannelInfo = "GetMotoStatusDownstreamChannelInfo"
	GetMotoStatusLog = "GetMotoStatusLog"
	GetMotoStatusSoftware = "GetMotoStatusSoftware"
	GetMotoStatusStartupSequence = "GetMotoStatusStartupSequence"
	GetMotoStatusUpstreamChannelInfo = "GetMotoStatusUpstreamChannelInfo"
)

// {
//     "GetMultipleHNAPsResponse": {
//         "GetMotoStatusStartupSequenceResponse": {
//             "MotoConnDSFreq": "663000000 Hz",
//             "MotoConnDSComment": "Locked",
//             "MotoConnConnectivityStatus": "OK",
//             "MotoConnConnectivityComment": "Operational",
//             "MotoConnBootStatus": "OK",
//             "MotoConnBootComment": "Operational",
//             "MotoConnConfigurationFileStatus": "OK",
//             "MotoConnConfigurationFileComment": "d11_m_mb8600_gigabit_c01.cm",
//             "MotoConnSecurityStatus": "Enabled",
//             "MotoConnSecurityComment": "BPI+",
//             "GetMotoStatusStartupSequenceResult": "OK"
//         },
//         "GetMotoStatusConnectionInfoResponse": {
//             "MotoConnSystemUpTime": "4 days 08h:57m:40s",
//             "MotoConnNetworkAccess": "Allowed",
//             "GetMotoStatusConnectionInfoResult": "OK"
//         },
//         "GetMotoStatusDownstreamChannelInfoResponse": {
//             "MotoConnDownstreamChannel": "1^Locked^QAM256^33^663.0^-9.3^38.8^42325^10482^|+|2^Locked^QAM256^5^483.0^-9.7^32.9^871509^86661^|+|3^Locked^QAM256^6^489.0^-10.3^33.0^1063003^100804^|+|4^Locked^QAM256^7^495.0^-10.0^37.0^876418^48477^|+|5^Locked^QAM256^8^507.0^-10.2^36.8^99458^9393^|+|6^Locked^QAM256^9^513.0^-9.8^37.6^57649^9245^|+|7^Locked^QAM256^10^519.0^-10.2^34.2^71405^10615^|+|8^Locked^Unknown^11^525.0^-9.9^ 0.0^0^0^|+|9^Locked^Unknown^12^531.0^-9.4^ 0.0^0^0^|+|10^Locked^QAM256^13^543.0^-10.4^ 0.0^57604742^6115113^|+|11^Locked^QAM256^14^549.0^-10.1^30.1^6865402^21651^|+|12^Locked^QAM256^15^555.0^-10.4^29.2^56081839^31637^|+|13^Locked^QAM256^16^561.0^-10.1^30.4^982473^22663^|+|14^Locked^Unknown^17^567.0^-9.8^ 0.0^0^0^|+|15^Locked^Unknown^18^573.0^-10.1^ 0.0^0^0^|+|16^Locked^Unknown^19^579.0^-9.1^ 0.0^25507^0^|+|17^Locked^QAM256^20^585.0^-10.1^36.5^46721^9090^|+|18^Locked^QAM256^21^591.0^-9.2^37.1^48459^9123^|+|19^Locked^QAM256^22^597.0^-9.3^35.3^55943^10332^|+|20^Locked^Unknown^23^603.0^-9.5^ 0.0^72126321^8229323^|+|21^Locked^QAM256^24^609.0^-8.6^31.2^365189731^1277974^|+|22^Locked^QAM256^25^615.0^-9.5^38.1^42659^8351^|+|23^Locked^QAM256^26^621.0^-8.8^38.3^40739^7864^|+|24^Locked^QAM256^27^627.0^-9.6^36.0^42685^8476^|+|25^Locked^QAM256^28^633.0^-9.6^34.4^43216^8440^|+|26^Locked^QAM256^29^639.0^-9.6^36.8^41998^8142^|+|27^Locked^QAM256^30^645.0^-9.4^35.7^41780^8288^|+|28^Locked^QAM256^31^651.0^-9.5^37.9^41818^8580^|+|29^Locked^QAM256^32^657.0^-9.8^36.9^44433^9910^|+|30^Locked^QAM256^34^669.0^-10.4^37.4^50186^20299^|+|31^Locked^QAM256^35^675.0^-9.5^37.3^53103^18787^|+|32^Locked^QAM256^36^681.0^-10.2^37.8^41924^11230^|+|33^Locked^OFDM PLC^159^722.0^-8.4^21.5^-1773898168^1086340^",
//             "GetMotoStatusDownstreamChannelInfoResult": "OK"
//         },
//         "GetMotoStatusUpstreamChannelInfoResponse": {
//             "MotoConnUpstreamChannel": "1^Locked^SC-QAM^1^5120^17.3^58.8^|+|2^Locked^SC-QAM^2^5120^23.7^58.8^|+|3^Locked^SC-QAM^3^5120^30.1^58.8^|+|4^Locked^SC-QAM^4^5120^36.5^58.8^",
//             "GetMotoStatusUpstreamChannelInfoResult": "OK"
//         },
//         "GetMotoLagStatusResponse": {
//             "MotoLagCurrentStatus": "0",
//             "GetMotoLagStatusResult": "OK"
//         },
//         "GetMultipleHNAPsResult": "OK"
//     }
// }
