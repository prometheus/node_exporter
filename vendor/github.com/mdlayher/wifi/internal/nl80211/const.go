// WARNING: This file has automatically been generated on Thu, 29 Dec 2016 08:53:57 EST.
// By https://git.io/cgogen. DO NOT EDIT.

package nl80211

const (
	// GenlName as defined in nl80211/nl80211.h:44
	GenlName = "nl80211"
	// MulticastGroupConfig as defined in nl80211/nl80211.h:46
	MulticastGroupConfig = "config"
	// MulticastGroupScan as defined in nl80211/nl80211.h:47
	MulticastGroupScan = "scan"
	// MulticastGroupReg as defined in nl80211/nl80211.h:48
	MulticastGroupReg = "regulatory"
	// MulticastGroupMlme as defined in nl80211/nl80211.h:49
	MulticastGroupMlme = "mlme"
	// MulticastGroupVendor as defined in nl80211/nl80211.h:50
	MulticastGroupVendor = "vendor"
	// MulticastGroupNan as defined in nl80211/nl80211.h:51
	MulticastGroupNan = "nan"
	// MulticastGroupTestmode as defined in nl80211/nl80211.h:52
	MulticastGroupTestmode = "testmode"
	// CmdGetMeshParams as defined in nl80211/nl80211.h:1080
	CmdGetMeshParams = CmdGetMeshConfig
	// CmdSetMeshParams as defined in nl80211/nl80211.h:1081
	CmdSetMeshParams = CmdSetMeshConfig
	// MeshSetupVendorPathSelIe as defined in nl80211/nl80211.h:1082
	MeshSetupVendorPathSelIe = MeshSetupIe
	// AttrScanGeneration as defined in nl80211/nl80211.h:2332
	AttrScanGeneration = AttrGeneration
	// AttrMeshParams as defined in nl80211/nl80211.h:2333
	AttrMeshParams = AttrMeshConfig
	// AttrIfaceSocketOwner as defined in nl80211/nl80211.h:2334
	AttrIfaceSocketOwner = AttrSocketOwner
	// MaxSuppRates as defined in nl80211/nl80211.h:2336
	MaxSuppRates = 32
	// MaxSuppHtRates as defined in nl80211/nl80211.h:2337
	MaxSuppHtRates = 77
	// MaxSuppRegRules as defined in nl80211/nl80211.h:2338
	MaxSuppRegRules = 64
	// TkipDataOffsetEncrKey as defined in nl80211/nl80211.h:2339
	TkipDataOffsetEncrKey = 0
	// TkipDataOffsetTxMicKey as defined in nl80211/nl80211.h:2340
	TkipDataOffsetTxMicKey = 16
	// TkipDataOffsetRxMicKey as defined in nl80211/nl80211.h:2341
	TkipDataOffsetRxMicKey = 24
	// HtCapabilityLen as defined in nl80211/nl80211.h:2342
	HtCapabilityLen = 26
	// VhtCapabilityLen as defined in nl80211/nl80211.h:2343
	VhtCapabilityLen = 12
	// MaxNrCipherSuites as defined in nl80211/nl80211.h:2345
	MaxNrCipherSuites = 5
	// MaxNrAkmSuites as defined in nl80211/nl80211.h:2346
	MaxNrAkmSuites = 2
	// MinRemainOnChannelTime as defined in nl80211/nl80211.h:2348
	MinRemainOnChannelTime = 10
	// ScanRssiTholdOff as defined in nl80211/nl80211.h:2351
	ScanRssiTholdOff = -300
	// CqmTxeMaxIntvl as defined in nl80211/nl80211.h:2353
	CqmTxeMaxIntvl = 1800
	// StaFlagMaxOldApi as defined in nl80211/nl80211.h:2457
	StaFlagMaxOldApi = StaFlagTdlsPeer
	// FrequencyAttrPassiveScan as defined in nl80211/nl80211.h:2858
	FrequencyAttrPassiveScan = FrequencyAttrNoIr
	// FrequencyAttrNoIbss as defined in nl80211/nl80211.h:2859
	FrequencyAttrNoIbss = FrequencyAttrNoIr
	// FrequencyAttrGoConcurrent as defined in nl80211/nl80211.h:2860
	FrequencyAttrGoConcurrent = FrequencyAttrIrConcurrent
	// AttrSchedScanMatchSsid as defined in nl80211/nl80211.h:3002
	AttrSchedScanMatchSsid = SchedScanMatchAttrSsid
	// RrfPassiveScan as defined in nl80211/nl80211.h:3044
	RrfPassiveScan = RrfNoIr
	// RrfNoIbss as defined in nl80211/nl80211.h:3045
	RrfNoIbss = RrfNoIr
	// RrfNoHt40 as defined in nl80211/nl80211.h:3046
	RrfNoHt40 = (RrfNoHt40minus | RrfNoHt40plus)
	// RrfGoConcurrent as defined in nl80211/nl80211.h:3048
	RrfGoConcurrent = RrfIrConcurrent
	// RrfNoIrAll as defined in nl80211/nl80211.h:3051
	RrfNoIrAll = (RrfNoIr | __RrfNoIbss)
	// SurveyInfoChannelTime as defined in nl80211/nl80211.h:3137
	SurveyInfoChannelTime = SurveyInfoTime
	// SurveyInfoChannelTimeBusy as defined in nl80211/nl80211.h:3138
	SurveyInfoChannelTimeBusy = SurveyInfoTimeBusy
	// SurveyInfoChannelTimeExtBusy as defined in nl80211/nl80211.h:3139
	SurveyInfoChannelTimeExtBusy = SurveyInfoTimeExtBusy
	// SurveyInfoChannelTimeRx as defined in nl80211/nl80211.h:3140
	SurveyInfoChannelTimeRx = SurveyInfoTimeRx
	// SurveyInfoChannelTimeTx as defined in nl80211/nl80211.h:3141
	SurveyInfoChannelTimeTx = SurveyInfoTimeTx
	// TxqAttrQueue as defined in nl80211/nl80211.h:3449
	TxqAttrQueue = TxqAttrAc
	// TxqQVo as defined in nl80211/nl80211.h:3450
	TxqQVo = AcVo
	// TxqQVi as defined in nl80211/nl80211.h:3451
	TxqQVi = AcVi
	// TxqQBe as defined in nl80211/nl80211.h:3452
	TxqQBe = AcBe
	// TxqQBk as defined in nl80211/nl80211.h:3453
	TxqQBk = AcBk
	// TxrateMcs as defined in nl80211/nl80211.h:3748
	TxrateMcs = TxrateHt
	// VhtNssMax as defined in nl80211/nl80211.h:3749
	VhtNssMax = 8
	// __WowlanPktpatInvalid as defined in nl80211/nl80211.h:3914
	__WowlanPktpatInvalid = __PktpatInvalid
	// WowlanPktpatMask as defined in nl80211/nl80211.h:3915
	WowlanPktpatMask = PktpatMask
	// WowlanPktpatPattern as defined in nl80211/nl80211.h:3916
	WowlanPktpatPattern = PktpatPattern
	// WowlanPktpatOffset as defined in nl80211/nl80211.h:3917
	WowlanPktpatOffset = PktpatOffset
	// Num_WowlanPktpat as defined in nl80211/nl80211.h:3918
	Num_WowlanPktpat = Num_Pktpat
	// Max_WowlanPktpat as defined in nl80211/nl80211.h:3919
	Max_WowlanPktpat = Max_Pktpat
	// WowlanPatternSupport as defined in nl80211/nl80211.h:3920
	WowlanPatternSupport = 0
	// KckLen as defined in nl80211/nl80211.h:4327
	KckLen = 16
	// KekLen as defined in nl80211/nl80211.h:4328
	KekLen = 16
	// ReplayCtrLen as defined in nl80211/nl80211.h:4329
	ReplayCtrLen = 8
	// CritProtoMaxDuration as defined in nl80211/nl80211.h:4781
	CritProtoMaxDuration = 5000
	// VendorIdIsLinux as defined in nl80211/nl80211.h:4799
	VendorIdIsLinux = 0x80000000
	// NanFuncServiceIdLen as defined in nl80211/nl80211.h:4960
	NanFuncServiceIdLen = 6
	// NanFuncServiceSpecInfoMaxLen as defined in nl80211/nl80211.h:4961
	NanFuncServiceSpecInfoMaxLen = 0xff
	// NanFuncSrfMaxLen as defined in nl80211/nl80211.h:4962
	NanFuncSrfMaxLen = 0xff
)

// nl80211Commands as declared in nl80211/nl80211.h:880
type nl80211Commands int32

// nl80211Commands enumeration from nl80211/nl80211.h:880
const (
	CmdUnspec                  = iota
	CmdGetWiphy                = 1
	CmdSetWiphy                = 2
	CmdNewWiphy                = 3
	CmdDelWiphy                = 4
	CmdGetInterface            = 5
	CmdSetInterface            = 6
	CmdNewInterface            = 7
	CmdDelInterface            = 8
	CmdGetKey                  = 9
	CmdSetKey                  = 10
	CmdNewKey                  = 11
	CmdDelKey                  = 12
	CmdGetBeacon               = 13
	CmdSetBeacon               = 14
	CmdStartAp                 = 15
	CmdNewBeacon               = CmdStartAp
	CmdStopAp                  = 16
	CmdDelBeacon               = CmdStopAp
	CmdGetStation              = 17
	CmdSetStation              = 18
	CmdNewStation              = 19
	CmdDelStation              = 20
	CmdGetMpath                = 21
	CmdSetMpath                = 22
	CmdNewMpath                = 23
	CmdDelMpath                = 24
	CmdSetBss                  = 25
	CmdSetReg                  = 26
	CmdReqSetReg               = 27
	CmdGetMeshConfig           = 28
	CmdSetMeshConfig           = 29
	CmdSetMgmtExtraIe          = 30
	CmdGetReg                  = 31
	CmdGetScan                 = 32
	CmdTriggerScan             = 33
	CmdNewScanResults          = 34
	CmdScanAborted             = 35
	CmdRegChange               = 36
	CmdAuthenticate            = 37
	CmdAssociate               = 38
	CmdDeauthenticate          = 39
	CmdDisassociate            = 40
	CmdMichaelMicFailure       = 41
	CmdRegBeaconHint           = 42
	CmdJoinIbss                = 43
	CmdLeaveIbss               = 44
	CmdTestmode                = 45
	CmdConnect                 = 46
	CmdRoam                    = 47
	CmdDisconnect              = 48
	CmdSetWiphyNetns           = 49
	CmdGetSurvey               = 50
	CmdNewSurveyResults        = 51
	CmdSetPmksa                = 52
	CmdDelPmksa                = 53
	CmdFlushPmksa              = 54
	CmdRemainOnChannel         = 55
	CmdCancelRemainOnChannel   = 56
	CmdSetTxBitrateMask        = 57
	CmdRegisterFrame           = 58
	CmdRegisterAction          = CmdRegisterFrame
	CmdFrame                   = 59
	CmdAction                  = CmdFrame
	CmdFrameTxStatus           = 60
	CmdActionTxStatus          = CmdFrameTxStatus
	CmdSetPowerSave            = 61
	CmdGetPowerSave            = 62
	CmdSetCqm                  = 63
	CmdNotifyCqm               = 64
	CmdSetChannel              = 65
	CmdSetWdsPeer              = 66
	CmdFrameWaitCancel         = 67
	CmdJoinMesh                = 68
	CmdLeaveMesh               = 69
	CmdUnprotDeauthenticate    = 70
	CmdUnprotDisassociate      = 71
	CmdNewPeerCandidate        = 72
	CmdGetWowlan               = 73
	CmdSetWowlan               = 74
	CmdStartSchedScan          = 75
	CmdStopSchedScan           = 76
	CmdSchedScanResults        = 77
	CmdSchedScanStopped        = 78
	CmdSetRekeyOffload         = 79
	CmdPmksaCandidate          = 80
	CmdTdlsOper                = 81
	CmdTdlsMgmt                = 82
	CmdUnexpectedFrame         = 83
	CmdProbeClient             = 84
	CmdRegisterBeacons         = 85
	CmdUnexpected4addrFrame    = 86
	CmdSetNoackMap             = 87
	CmdChSwitchNotify          = 88
	CmdStartP2pDevice          = 89
	CmdStopP2pDevice           = 90
	CmdConnFailed              = 91
	CmdSetMcastRate            = 92
	CmdSetMacAcl               = 93
	CmdRadarDetect             = 94
	CmdGetProtocolFeatures     = 95
	CmdUpdateFtIes             = 96
	CmdFtEvent                 = 97
	CmdCritProtocolStart       = 98
	CmdCritProtocolStop        = 99
	CmdGetCoalesce             = 100
	CmdSetCoalesce             = 101
	CmdChannelSwitch           = 102
	CmdVendor                  = 103
	CmdSetQosMap               = 104
	CmdAddTxTs                 = 105
	CmdDelTxTs                 = 106
	CmdGetMpp                  = 107
	CmdJoinOcb                 = 108
	CmdLeaveOcb                = 109
	CmdChSwitchStartedNotify   = 110
	CmdTdlsChannelSwitch       = 111
	CmdTdlsCancelChannelSwitch = 112
	CmdWiphyRegChange          = 113
	CmdAbortScan               = 114
	CmdStartNan                = 115
	CmdStopNan                 = 116
	CmdAddNanFunction          = 117
	CmdDelNanFunction          = 118
	CmdChangeNanConfig         = 119
	CmdNanMatch                = 120
	__CmdAfterLast             = 121
	CmdMax                     = __CmdAfterLast - 1
)

// nl80211Attrs as declared in nl80211/nl80211.h:1929
type nl80211Attrs int32

// nl80211Attrs enumeration from nl80211/nl80211.h:1929
const (
	AttrUnspec                       = iota
	AttrWiphy                        = 1
	AttrWiphyName                    = 2
	AttrIfindex                      = 3
	AttrIfname                       = 4
	AttrIftype                       = 5
	AttrMac                          = 6
	AttrKeyData                      = 7
	AttrKeyIdx                       = 8
	AttrKeyCipher                    = 9
	AttrKeySeq                       = 10
	AttrKeyDefault                   = 11
	AttrBeaconInterval               = 12
	AttrDtimPeriod                   = 13
	AttrBeaconHead                   = 14
	AttrBeaconTail                   = 15
	AttrStaAid                       = 16
	AttrStaFlags                     = 17
	AttrStaListenInterval            = 18
	AttrStaSupportedRates            = 19
	AttrStaVlan                      = 20
	AttrStaInfo                      = 21
	AttrWiphyBands                   = 22
	AttrMntrFlags                    = 23
	AttrMeshId                       = 24
	AttrStaPlinkAction               = 25
	AttrMpathNextHop                 = 26
	AttrMpathInfo                    = 27
	AttrBssCtsProt                   = 28
	AttrBssShortPreamble             = 29
	AttrBssShortSlotTime             = 30
	AttrHtCapability                 = 31
	AttrSupportedIftypes             = 32
	AttrRegAlpha2                    = 33
	AttrRegRules                     = 34
	AttrMeshConfig                   = 35
	AttrBssBasicRates                = 36
	AttrWiphyTxqParams               = 37
	AttrWiphyFreq                    = 38
	AttrWiphyChannelType             = 39
	AttrKeyDefaultMgmt               = 40
	AttrMgmtSubtype                  = 41
	AttrIe                           = 42
	AttrMaxNumScanSsids              = 43
	AttrScanFrequencies              = 44
	AttrScanSsids                    = 45
	AttrGeneration                   = 46
	AttrBss                          = 47
	AttrRegInitiator                 = 48
	AttrRegType                      = 49
	AttrSupportedCommands            = 50
	AttrFrame                        = 51
	AttrSsid                         = 52
	AttrAuthType                     = 53
	AttrReasonCode                   = 54
	AttrKeyType                      = 55
	AttrMaxScanIeLen                 = 56
	AttrCipherSuites                 = 57
	AttrFreqBefore                   = 58
	AttrFreqAfter                    = 59
	AttrFreqFixed                    = 60
	AttrWiphyRetryShort              = 61
	AttrWiphyRetryLong               = 62
	AttrWiphyFragThreshold           = 63
	AttrWiphyRtsThreshold            = 64
	AttrTimedOut                     = 65
	AttrUseMfp                       = 66
	AttrStaFlags2                    = 67
	AttrControlPort                  = 68
	AttrTestdata                     = 69
	AttrPrivacy                      = 70
	AttrDisconnectedByAp             = 71
	AttrStatusCode                   = 72
	AttrCipherSuitesPairwise         = 73
	AttrCipherSuiteGroup             = 74
	AttrWpaVersions                  = 75
	AttrAkmSuites                    = 76
	AttrReqIe                        = 77
	AttrRespIe                       = 78
	AttrPrevBssid                    = 79
	AttrKey                          = 80
	AttrKeys                         = 81
	AttrPid                          = 82
	Attr4addr                        = 83
	AttrSurveyInfo                   = 84
	AttrPmkid                        = 85
	AttrMaxNumPmkids                 = 86
	AttrDuration                     = 87
	AttrCookie                       = 88
	AttrWiphyCoverageClass           = 89
	AttrTxRates                      = 90
	AttrFrameMatch                   = 91
	AttrAck                          = 92
	AttrPsState                      = 93
	AttrCqm                          = 94
	AttrLocalStateChange             = 95
	AttrApIsolate                    = 96
	AttrWiphyTxPowerSetting          = 97
	AttrWiphyTxPowerLevel            = 98
	AttrTxFrameTypes                 = 99
	AttrRxFrameTypes                 = 100
	AttrFrameType                    = 101
	AttrControlPortEthertype         = 102
	AttrControlPortNoEncrypt         = 103
	AttrSupportIbssRsn               = 104
	AttrWiphyAntennaTx               = 105
	AttrWiphyAntennaRx               = 106
	AttrMcastRate                    = 107
	AttrOffchannelTxOk               = 108
	AttrBssHtOpmode                  = 109
	AttrKeyDefaultTypes              = 110
	AttrMaxRemainOnChannelDuration   = 111
	AttrMeshSetup                    = 112
	AttrWiphyAntennaAvailTx          = 113
	AttrWiphyAntennaAvailRx          = 114
	AttrSupportMeshAuth              = 115
	AttrStaPlinkState                = 116
	AttrWowlanTriggers               = 117
	AttrWowlanTriggersSupported      = 118
	AttrSchedScanInterval            = 119
	AttrInterfaceCombinations        = 120
	AttrSoftwareIftypes              = 121
	AttrRekeyData                    = 122
	AttrMaxNumSchedScanSsids         = 123
	AttrMaxSchedScanIeLen            = 124
	AttrScanSuppRates                = 125
	AttrHiddenSsid                   = 126
	AttrIeProbeResp                  = 127
	AttrIeAssocResp                  = 128
	AttrStaWme                       = 129
	AttrSupportApUapsd               = 130
	AttrRoamSupport                  = 131
	AttrSchedScanMatch               = 132
	AttrMaxMatchSets                 = 133
	AttrPmksaCandidate               = 134
	AttrTxNoCckRate                  = 135
	AttrTdlsAction                   = 136
	AttrTdlsDialogToken              = 137
	AttrTdlsOperation                = 138
	AttrTdlsSupport                  = 139
	AttrTdlsExternalSetup            = 140
	AttrDeviceApSme                  = 141
	AttrDontWaitForAck               = 142
	AttrFeatureFlags                 = 143
	AttrProbeRespOffload             = 144
	AttrProbeResp                    = 145
	AttrDfsRegion                    = 146
	AttrDisableHt                    = 147
	AttrHtCapabilityMask             = 148
	AttrNoackMap                     = 149
	AttrInactivityTimeout            = 150
	AttrRxSignalDbm                  = 151
	AttrBgScanPeriod                 = 152
	AttrWdev                         = 153
	AttrUserRegHintType              = 154
	AttrConnFailedReason             = 155
	AttrSaeData                      = 156
	AttrVhtCapability                = 157
	AttrScanFlags                    = 158
	AttrChannelWidth                 = 159
	AttrCenterFreq1                  = 160
	AttrCenterFreq2                  = 161
	AttrP2pCtwindow                  = 162
	AttrP2pOppps                     = 163
	AttrLocalMeshPowerMode           = 164
	AttrAclPolicy                    = 165
	AttrMacAddrs                     = 166
	AttrMacAclMax                    = 167
	AttrRadarEvent                   = 168
	AttrExtCapa                      = 169
	AttrExtCapaMask                  = 170
	AttrStaCapability                = 171
	AttrStaExtCapability             = 172
	AttrProtocolFeatures             = 173
	AttrSplitWiphyDump               = 174
	AttrDisableVht                   = 175
	AttrVhtCapabilityMask            = 176
	AttrMdid                         = 177
	AttrIeRic                        = 178
	AttrCritProtId                   = 179
	AttrMaxCritProtDuration          = 180
	AttrPeerAid                      = 181
	AttrCoalesceRule                 = 182
	AttrChSwitchCount                = 183
	AttrChSwitchBlockTx              = 184
	AttrCsaIes                       = 185
	AttrCsaCOffBeacon                = 186
	AttrCsaCOffPresp                 = 187
	AttrRxmgmtFlags                  = 188
	AttrStaSupportedChannels         = 189
	AttrStaSupportedOperClasses      = 190
	AttrHandleDfs                    = 191
	AttrSupport5Mhz                  = 192
	AttrSupport10Mhz                 = 193
	AttrOpmodeNotif                  = 194
	AttrVendorId                     = 195
	AttrVendorSubcmd                 = 196
	AttrVendorData                   = 197
	AttrVendorEvents                 = 198
	AttrQosMap                       = 199
	AttrMacHint                      = 200
	AttrWiphyFreqHint                = 201
	AttrMaxApAssocSta                = 202
	AttrTdlsPeerCapability           = 203
	AttrSocketOwner                  = 204
	AttrCsaCOffsetsTx                = 205
	AttrMaxCsaCounters               = 206
	AttrTdlsInitiator                = 207
	AttrUseRrm                       = 208
	AttrWiphyDynAck                  = 209
	AttrTsid                         = 210
	AttrUserPrio                     = 211
	AttrAdmittedTime                 = 212
	AttrSmpsMode                     = 213
	AttrOperClass                    = 214
	AttrMacMask                      = 215
	AttrWiphySelfManagedReg          = 216
	AttrExtFeatures                  = 217
	AttrSurveyRadioStats             = 218
	AttrNetnsFd                      = 219
	AttrSchedScanDelay               = 220
	AttrRegIndoor                    = 221
	AttrMaxNumSchedScanPlans         = 222
	AttrMaxScanPlanInterval          = 223
	AttrMaxScanPlanIterations        = 224
	AttrSchedScanPlans               = 225
	AttrPbss                         = 226
	AttrBssSelect                    = 227
	AttrStaSupportP2pPs              = 228
	AttrPad                          = 229
	AttrIftypeExtCapa                = 230
	AttrMuMimoGroupData              = 231
	AttrMuMimoFollowMacAddr          = 232
	AttrScanStartTimeTsf             = 233
	AttrScanStartTimeTsfBssid        = 234
	AttrMeasurementDuration          = 235
	AttrMeasurementDurationMandatory = 236
	AttrMeshPeerAid                  = 237
	AttrNanMasterPref                = 238
	AttrNanDual                      = 239
	AttrNanFunc                      = 240
	AttrNanMatch                     = 241
	__AttrAfterLast                  = 242
	Num_Attr                         = __AttrAfterLast
	AttrMax                          = __AttrAfterLast - 1
)

// nl80211Iftype as declared in nl80211/nl80211.h:2384
type nl80211Iftype int32

// nl80211Iftype enumeration from nl80211/nl80211.h:2384
const (
	IftypeUnspecified = iota
	IftypeAdhoc       = 1
	IftypeStation     = 2
	IftypeAp          = 3
	IftypeApVlan      = 4
	IftypeWds         = 5
	IftypeMonitor     = 6
	IftypeMeshPoint   = 7
	IftypeP2pClient   = 8
	IftypeP2pGo       = 9
	IftypeP2pDevice   = 10
	IftypeOcb         = 11
	IftypeNan         = 12
	Num_Iftypes       = 13
	IftypeMax         = Num_Iftypes - 1
)

// nl80211StaFlags as declared in nl80211/nl80211.h:2428
type nl80211StaFlags int32

// nl80211StaFlags enumeration from nl80211/nl80211.h:2428
const (
	__StaFlagInvalid     = iota
	StaFlagAuthorized    = 1
	StaFlagShortPreamble = 2
	StaFlagWme           = 3
	StaFlagMfp           = 4
	StaFlagAuthenticated = 5
	StaFlagTdlsPeer      = 6
	StaFlagAssociated    = 7
	__StaFlagAfterLast   = 8
	StaFlagMax           = __StaFlagAfterLast - 1
)

// nl80211StaP2pPsStatus as declared in nl80211/nl80211.h:2450
type nl80211StaP2pPsStatus int32

// nl80211StaP2pPsStatus enumeration from nl80211/nl80211.h:2450
const (
	P2pPsUnsupported = iota
	P2pPsSupported   = 1
	Num_P2pPsStatus  = 2
)

// nl80211RateInfo as declared in nl80211/nl80211.h:2505
type nl80211RateInfo int32

// nl80211RateInfo enumeration from nl80211/nl80211.h:2505
const (
	__RateInfoInvalid     = iota
	RateInfoBitrate       = 1
	RateInfoMcs           = 2
	RateInfo40MhzWidth    = 3
	RateInfoShortGi       = 4
	RateInfoBitrate32     = 5
	RateInfoVhtMcs        = 6
	RateInfoVhtNss        = 7
	RateInfo80MhzWidth    = 8
	RateInfo80p80MhzWidth = 9
	RateInfo160MhzWidth   = 10
	RateInfo10MhzWidth    = 11
	RateInfo5MhzWidth     = 12
	__RateInfoAfterLast   = 13
	RateInfoMax           = __RateInfoAfterLast - 1
)

// nl80211StaBssParam as declared in nl80211/nl80211.h:2542
type nl80211StaBssParam int32

// nl80211StaBssParam enumeration from nl80211/nl80211.h:2542
const (
	__StaBssParamInvalid      = iota
	StaBssParamCtsProt        = 1
	StaBssParamShortPreamble  = 2
	StaBssParamShortSlotTime  = 3
	StaBssParamDtimPeriod     = 4
	StaBssParamBeaconInterval = 5
	__StaBssParamAfterLast    = 6
	StaBssParamMax            = __StaBssParamAfterLast - 1
)

// nl80211StaInfo as declared in nl80211/nl80211.h:2620
type nl80211StaInfo int32

// nl80211StaInfo enumeration from nl80211/nl80211.h:2620
const (
	__StaInfoInvalid          = iota
	StaInfoInactiveTime       = 1
	StaInfoRxBytes            = 2
	StaInfoTxBytes            = 3
	StaInfoLlid               = 4
	StaInfoPlid               = 5
	StaInfoPlinkState         = 6
	StaInfoSignal             = 7
	StaInfoTxBitrate          = 8
	StaInfoRxPackets          = 9
	StaInfoTxPackets          = 10
	StaInfoTxRetries          = 11
	StaInfoTxFailed           = 12
	StaInfoSignalAvg          = 13
	StaInfoRxBitrate          = 14
	StaInfoBssParam           = 15
	StaInfoConnectedTime      = 16
	StaInfoStaFlags           = 17
	StaInfoBeaconLoss         = 18
	StaInfoTOffset            = 19
	StaInfoLocalPm            = 20
	StaInfoPeerPm             = 21
	StaInfoNonpeerPm          = 22
	StaInfoRxBytes64          = 23
	StaInfoTxBytes64          = 24
	StaInfoChainSignal        = 25
	StaInfoChainSignalAvg     = 26
	StaInfoExpectedThroughput = 27
	StaInfoRxDropMisc         = 28
	StaInfoBeaconRx           = 29
	StaInfoBeaconSignalAvg    = 30
	StaInfoTidStats           = 31
	StaInfoRxDuration         = 32
	StaInfoPad                = 33
	__StaInfoAfterLast        = 34
	StaInfoMax                = __StaInfoAfterLast - 1
)

// nl80211TidStats as declared in nl80211/nl80211.h:2675
type nl80211TidStats int32

// nl80211TidStats enumeration from nl80211/nl80211.h:2675
const (
	__TidStatsInvalid     = iota
	TidStatsRxMsdu        = 1
	TidStatsTxMsdu        = 2
	TidStatsTxMsduRetries = 3
	TidStatsTxMsduFailed  = 4
	TidStatsPad           = 5
	Num_TidStats          = 6
	TidStatsMax           = Num_TidStats - 1
)

// nl80211MpathFlags as declared in nl80211/nl80211.h:2697
type nl80211MpathFlags int32

// nl80211MpathFlags enumeration from nl80211/nl80211.h:2697
const (
	MpathFlagActive    = 1 << 0
	MpathFlagResolving = 1 << 1
	MpathFlagSnValid   = 1 << 2
	MpathFlagFixed     = 1 << 3
	MpathFlagResolved  = 1 << 4
)

// nl80211MpathInfo as declared in nl80211/nl80211.h:2724
type nl80211MpathInfo int32

// nl80211MpathInfo enumeration from nl80211/nl80211.h:2724
const (
	__MpathInfoInvalid        = iota
	MpathInfoFrameQlen        = 1
	MpathInfoSn               = 2
	MpathInfoMetric           = 3
	MpathInfoExptime          = 4
	MpathInfoFlags            = 5
	MpathInfoDiscoveryTimeout = 6
	MpathInfoDiscoveryRetries = 7
	__MpathInfoAfterLast      = 8
	MpathInfoMax              = __MpathInfoAfterLast - 1
)

// nl80211BandAttr as declared in nl80211/nl80211.h:2757
type nl80211BandAttr int32

// nl80211BandAttr enumeration from nl80211/nl80211.h:2757
const (
	__BandAttrInvalid      = iota
	BandAttrFreqs          = 1
	BandAttrRates          = 2
	BandAttrHtMcsSet       = 3
	BandAttrHtCapa         = 4
	BandAttrHtAmpduFactor  = 5
	BandAttrHtAmpduDensity = 6
	BandAttrVhtMcsSet      = 7
	BandAttrVhtCapa        = 8
	__BandAttrAfterLast    = 9
	BandAttrMax            = __BandAttrAfterLast - 1
)

// nl80211FrequencyAttr as declared in nl80211/nl80211.h:2833
type nl80211FrequencyAttr int32

// nl80211FrequencyAttr enumeration from nl80211/nl80211.h:2833
const (
	__FrequencyAttrInvalid    = iota
	FrequencyAttrFreq         = 1
	FrequencyAttrDisabled     = 2
	FrequencyAttrNoIr         = 3
	__FrequencyAttrNoIbss     = 4
	FrequencyAttrRadar        = 5
	FrequencyAttrMaxTxPower   = 6
	FrequencyAttrDfsState     = 7
	FrequencyAttrDfsTime      = 8
	FrequencyAttrNoHt40Minus  = 9
	FrequencyAttrNoHt40Plus   = 10
	FrequencyAttrNo80mhz      = 11
	FrequencyAttrNo160mhz     = 12
	FrequencyAttrDfsCacTime   = 13
	FrequencyAttrIndoorOnly   = 14
	FrequencyAttrIrConcurrent = 15
	FrequencyAttrNo20mhz      = 16
	FrequencyAttrNo10mhz      = 17
	__FrequencyAttrAfterLast  = 18
	FrequencyAttrMax          = __FrequencyAttrAfterLast - 1
)

// nl80211BitrateAttr as declared in nl80211/nl80211.h:2873
type nl80211BitrateAttr int32

// nl80211BitrateAttr enumeration from nl80211/nl80211.h:2873
const (
	__BitrateAttrInvalid         = iota
	BitrateAttrRate              = 1
	BitrateAttr2ghzShortpreamble = 2
	__BitrateAttrAfterLast       = 3
	BitrateAttrMax               = __BitrateAttrAfterLast - 1
)

// nl80211RegInitiator as declared in nl80211/nl80211.h:2899
type nl80211RegInitiator int32

// nl80211RegInitiator enumeration from nl80211/nl80211.h:2899
const (
	RegdomSetByCore      = iota
	RegdomSetByUser      = 1
	RegdomSetByDriver    = 2
	RegdomSetByCountryIe = 3
)

// nl80211RegType as declared in nl80211/nl80211.h:2922
type nl80211RegType int32

// nl80211RegType enumeration from nl80211/nl80211.h:2922
const (
	RegdomTypeCountry      = iota
	RegdomTypeWorld        = 1
	RegdomTypeCustomWorld  = 2
	RegdomTypeIntersection = 3
)

// nl80211RegRuleAttr as declared in nl80211/nl80211.h:2954
type nl80211RegRuleAttr int32

// nl80211RegRuleAttr enumeration from nl80211/nl80211.h:2954
const (
	__RegRuleAttrInvalid    = iota
	AttrRegRuleFlags        = 1
	AttrFreqRangeStart      = 2
	AttrFreqRangeEnd        = 3
	AttrFreqRangeMaxBw      = 4
	AttrPowerRuleMaxAntGain = 5
	AttrPowerRuleMaxEirp    = 6
	AttrDfsCacTime          = 7
	__RegRuleAttrAfterLast  = 8
	RegRuleAttrMax          = __RegRuleAttrAfterLast - 1
)

// nl80211SchedScanMatchAttr as declared in nl80211/nl80211.h:2989
type nl80211SchedScanMatchAttr int32

// nl80211SchedScanMatchAttr enumeration from nl80211/nl80211.h:2989
const (
	__SchedScanMatchAttrInvalid   = iota
	SchedScanMatchAttrSsid        = 1
	SchedScanMatchAttrRssi        = 2
	__SchedScanMatchAttrAfterLast = 3
	SchedScanMatchAttrMax         = __SchedScanMatchAttrAfterLast - 1
)

// nl80211RegRuleFlags as declared in nl80211/nl80211.h:3026
type nl80211RegRuleFlags int32

// nl80211RegRuleFlags enumeration from nl80211/nl80211.h:3026
const (
	RrfNoOfdm       = 1 << 0
	RrfNoCck        = 1 << 1
	RrfNoIndoor     = 1 << 2
	RrfNoOutdoor    = 1 << 3
	RrfDfs          = 1 << 4
	RrfPtpOnly      = 1 << 5
	RrfPtmpOnly     = 1 << 6
	RrfNoIr         = 1 << 7
	__RrfNoIbss     = 1 << 8
	RrfAutoBw       = 1 << 11
	RrfIrConcurrent = 1 << 12
	RrfNoHt40minus  = 1 << 13
	RrfNoHt40plus   = 1 << 14
	RrfNo80mhz      = 1 << 15
	RrfNo160mhz     = 1 << 16
)

// nl80211DfsRegions as declared in nl80211/nl80211.h:3061
type nl80211DfsRegions int32

// nl80211DfsRegions enumeration from nl80211/nl80211.h:3061
const (
	DfsUnset = iota
	DfsFcc   = 1
	DfsEtsi  = 2
	DfsJp    = 3
)

// nl80211UserRegHintType as declared in nl80211/nl80211.h:3085
type nl80211UserRegHintType int32

// nl80211UserRegHintType enumeration from nl80211/nl80211.h:3085
const (
	UserRegHintUser     = iota
	UserRegHintCellBase = 1
	UserRegHintIndoor   = 2
)

// nl80211SurveyInfo as declared in nl80211/nl80211.h:3118
type nl80211SurveyInfo int32

// nl80211SurveyInfo enumeration from nl80211/nl80211.h:3118
const (
	__SurveyInfoInvalid   = iota
	SurveyInfoFrequency   = 1
	SurveyInfoNoise       = 2
	SurveyInfoInUse       = 3
	SurveyInfoTime        = 4
	SurveyInfoTimeBusy    = 5
	SurveyInfoTimeExtBusy = 6
	SurveyInfoTimeRx      = 7
	SurveyInfoTimeTx      = 8
	SurveyInfoTimeScan    = 9
	SurveyInfoPad         = 10
	__SurveyInfoAfterLast = 11
	SurveyInfoMax         = __SurveyInfoAfterLast - 1
)

// nl80211MntrFlags as declared in nl80211/nl80211.h:3162
type nl80211MntrFlags int32

// nl80211MntrFlags enumeration from nl80211/nl80211.h:3162
const (
	__MntrFlagInvalid   = iota
	MntrFlagFcsfail     = 1
	MntrFlagPlcpfail    = 2
	MntrFlagControl     = 3
	MntrFlagOtherBss    = 4
	MntrFlagCookFrames  = 5
	MntrFlagActive      = 6
	__MntrFlagAfterLast = 7
	MntrFlagMax         = __MntrFlagAfterLast - 1
)

// nl80211MeshPowerMode as declared in nl80211/nl80211.h:3194
type nl80211MeshPowerMode int32

// nl80211MeshPowerMode enumeration from nl80211/nl80211.h:3194
const (
	MeshPowerUnknown     = iota
	MeshPowerActive      = 1
	MeshPowerLightSleep  = 2
	MeshPowerDeepSleep   = 3
	__MeshPowerAfterLast = 4
	MeshPowerMax         = __MeshPowerAfterLast - 1
)

// nl80211MeshconfParams as declared in nl80211/nl80211.h:3312
type nl80211MeshconfParams int32

// nl80211MeshconfParams enumeration from nl80211/nl80211.h:3312
const (
	__MeshconfInvalid                = iota
	MeshconfRetryTimeout             = 1
	MeshconfConfirmTimeout           = 2
	MeshconfHoldingTimeout           = 3
	MeshconfMaxPeerLinks             = 4
	MeshconfMaxRetries               = 5
	MeshconfTtl                      = 6
	MeshconfAutoOpenPlinks           = 7
	MeshconfHwmpMaxPreqRetries       = 8
	MeshconfPathRefreshTime          = 9
	MeshconfMinDiscoveryTimeout      = 10
	MeshconfHwmpActivePathTimeout    = 11
	MeshconfHwmpPreqMinInterval      = 12
	MeshconfHwmpNetDiamTrvsTime      = 13
	MeshconfHwmpRootmode             = 14
	MeshconfElementTtl               = 15
	MeshconfHwmpRannInterval         = 16
	MeshconfGateAnnouncements        = 17
	MeshconfHwmpPerrMinInterval      = 18
	MeshconfForwarding               = 19
	MeshconfRssiThreshold            = 20
	MeshconfSyncOffsetMaxNeighbor    = 21
	MeshconfHtOpmode                 = 22
	MeshconfHwmpPathToRootTimeout    = 23
	MeshconfHwmpRootInterval         = 24
	MeshconfHwmpConfirmationInterval = 25
	MeshconfPowerMode                = 26
	MeshconfAwakeWindow              = 27
	MeshconfPlinkTimeout             = 28
	__MeshconfAttrAfterLast          = 29
	MeshconfAttrMax                  = __MeshconfAttrAfterLast - 1
)

// nl80211MeshSetupParams as declared in nl80211/nl80211.h:3397
type nl80211MeshSetupParams int32

// nl80211MeshSetupParams enumeration from nl80211/nl80211.h:3397
const (
	__MeshSetupInvalid           = iota
	MeshSetupEnableVendorPathSel = 1
	MeshSetupEnableVendorMetric  = 2
	MeshSetupIe                  = 3
	MeshSetupUserspaceAuth       = 4
	MeshSetupUserspaceAmpe       = 5
	MeshSetupEnableVendorSync    = 6
	MeshSetupUserspaceMpm        = 7
	MeshSetupAuthProtocol        = 8
	__MeshSetupAttrAfterLast     = 9
	MeshSetupAttrMax             = __MeshSetupAttrAfterLast - 1
)

// nl80211TxqAttr as declared in nl80211/nl80211.h:3427
type nl80211TxqAttr int32

// nl80211TxqAttr enumeration from nl80211/nl80211.h:3427
const (
	__TxqAttrInvalid   = iota
	TxqAttrAc          = 1
	TxqAttrTxop        = 2
	TxqAttrCwmin       = 3
	TxqAttrCwmax       = 4
	TxqAttrAifs        = 5
	__TxqAttrAfterLast = 6
	TxqAttrMax         = __TxqAttrAfterLast - 1
)

// nl80211Ac as declared in nl80211/nl80211.h:3440
type nl80211Ac int32

// nl80211Ac enumeration from nl80211/nl80211.h:3440
const (
	AcVo   = iota
	AcVi   = 1
	AcBe   = 2
	AcBk   = 3
	NumAcs = 4
)

// nl80211ChannelType as declared in nl80211/nl80211.h:3464
type nl80211ChannelType int32

// nl80211ChannelType enumeration from nl80211/nl80211.h:3464
const (
	ChanNoHt      = iota
	ChanHt20      = 1
	ChanHt40minus = 2
	ChanHt40plus  = 3
)

// nl80211ChanWidth as declared in nl80211/nl80211.h:3490
type nl80211ChanWidth int32

// nl80211ChanWidth enumeration from nl80211/nl80211.h:3490
const (
	ChanWidth20Noht = iota
	ChanWidth20     = 1
	ChanWidth40     = 2
	ChanWidth80     = 3
	ChanWidth80p80  = 4
	ChanWidth160    = 5
	ChanWidth5      = 6
	ChanWidth10     = 7
)

// nl80211BssScanWidth as declared in nl80211/nl80211.h:3510
type nl80211BssScanWidth int32

// nl80211BssScanWidth enumeration from nl80211/nl80211.h:3510
const (
	BssChanWidth20 = iota
	BssChanWidth10 = 1
	BssChanWidth5  = 2
)

// nl80211Bss as declared in nl80211/nl80211.h:3565
type nl80211Bss int32

// nl80211Bss enumeration from nl80211/nl80211.h:3565
const (
	__BssInvalid           = iota
	BssBssid               = 1
	BssFrequency           = 2
	BssTsf                 = 3
	BssBeaconInterval      = 4
	BssCapability          = 5
	BssInformationElements = 6
	BssSignalMbm           = 7
	BssSignalUnspec        = 8
	BssStatus              = 9
	BssSeenMsAgo           = 10
	BssBeaconIes           = 11
	BssChanWidth           = 12
	BssBeaconTsf           = 13
	BssPrespData           = 14
	BssLastSeenBoottime    = 15
	BssPad                 = 16
	BssParentTsf           = 17
	BssParentBssid         = 18
	__BssAfterLast         = 19
	BssMax                 = __BssAfterLast - 1
)

// nl80211BssStatus as declared in nl80211/nl80211.h:3603
type nl80211BssStatus int32

// nl80211BssStatus enumeration from nl80211/nl80211.h:3603
const (
	BssStatusAuthenticated = iota
	BssStatusAssociated    = 1
	BssStatusIbssJoined    = 2
)

// nl80211AuthType as declared in nl80211/nl80211.h:3623
type nl80211AuthType int32

// nl80211AuthType enumeration from nl80211/nl80211.h:3623
const (
	AuthtypeOpenSystem = iota
	AuthtypeSharedKey  = 1
	AuthtypeFt         = 2
	AuthtypeNetworkEap = 3
	AuthtypeSae        = 4
	__AuthtypeNum      = 5
	AuthtypeMax        = __AuthtypeNum - 1
	AuthtypeAutomatic  = 5
)

// nl80211KeyType as declared in nl80211/nl80211.h:3643
type nl80211KeyType int32

// nl80211KeyType enumeration from nl80211/nl80211.h:3643
const (
	KeytypeGroup    = iota
	KeytypePairwise = 1
	KeytypePeerkey  = 2
	Num_Keytypes    = 3
)

// nl80211Mfp as declared in nl80211/nl80211.h:3656
type nl80211Mfp int32

// nl80211Mfp enumeration from nl80211/nl80211.h:3656
const (
	MfpNo       = iota
	MfpRequired = 1
)

// nl80211WpaVersions as declared in nl80211/nl80211.h:3661
type nl80211WpaVersions int32

// nl80211WpaVersions enumeration from nl80211/nl80211.h:3661
const (
	WpaVersion1 = 1 << 0
	WpaVersion2 = 1 << 1
)

// nl80211KeyDefaultTypes as declared in nl80211/nl80211.h:3675
type nl80211KeyDefaultTypes int32

// nl80211KeyDefaultTypes enumeration from nl80211/nl80211.h:3675
const (
	__KeyDefaultTypeInvalid = iota
	KeyDefaultTypeUnicast   = 1
	KeyDefaultTypeMulticast = 2
	Num_KeyDefaultTypes     = 3
)

// nl80211KeyAttributes as declared in nl80211/nl80211.h:3705
type nl80211KeyAttributes int32

// nl80211KeyAttributes enumeration from nl80211/nl80211.h:3705
const (
	__KeyInvalid    = iota
	KeyData         = 1
	KeyIdx          = 2
	KeyCipher       = 3
	KeySeq          = 4
	KeyDefault      = 5
	KeyDefaultMgmt  = 6
	KeyType         = 7
	KeyDefaultTypes = 8
	__KeyAfterLast  = 9
	KeyMax          = __KeyAfterLast - 1
)

// nl80211TxRateAttributes as declared in nl80211/nl80211.h:3736
type nl80211TxRateAttributes int32

// nl80211TxRateAttributes enumeration from nl80211/nl80211.h:3736
const (
	__TxrateInvalid   = iota
	TxrateLegacy      = 1
	TxrateHt          = 2
	TxrateVht         = 3
	TxrateGi          = 4
	__TxrateAfterLast = 5
	TxrateMax         = __TxrateAfterLast - 1
)

// nl80211TxrateGi as declared in nl80211/nl80211.h:3759
type nl80211TxrateGi int32

// nl80211TxrateGi enumeration from nl80211/nl80211.h:3759
const (
	TxrateDefaultGi = iota
	TxrateForceSgi  = 1
	TxrateForceLgi  = 2
)

// nl80211Band as declared in nl80211/nl80211.h:3773
type nl80211Band int32

// nl80211Band enumeration from nl80211/nl80211.h:3773
const (
	Band2ghz  = iota
	Band5ghz  = 1
	Band60ghz = 2
	Num_Bands = 3
)

// nl80211PsState as declared in nl80211/nl80211.h:3786
type nl80211PsState int32

// nl80211PsState enumeration from nl80211/nl80211.h:3786
const (
	PsDisabled = iota
	PsEnabled  = 1
)

// nl80211AttrCqm as declared in nl80211/nl80211.h:3819
type nl80211AttrCqm int32

// nl80211AttrCqm enumeration from nl80211/nl80211.h:3819
const (
	__AttrCqmInvalid          = iota
	AttrCqmRssiThold          = 1
	AttrCqmRssiHyst           = 2
	AttrCqmRssiThresholdEvent = 3
	AttrCqmPktLossEvent       = 4
	AttrCqmTxeRate            = 5
	AttrCqmTxePkts            = 6
	AttrCqmTxeIntvl           = 7
	AttrCqmBeaconLossEvent    = 8
	__AttrCqmAfterLast        = 9
	AttrCqmMax                = __AttrCqmAfterLast - 1
)

// nl80211CqmRssiThresholdEvent as declared in nl80211/nl80211.h:3843
type nl80211CqmRssiThresholdEvent int32

// nl80211CqmRssiThresholdEvent enumeration from nl80211/nl80211.h:3843
const (
	CqmRssiThresholdEventLow  = iota
	CqmRssiThresholdEventHigh = 1
	CqmRssiBeaconLossEvent    = 2
)

// nl80211TxPowerSetting as declared in nl80211/nl80211.h:3856
type nl80211TxPowerSetting int32

// nl80211TxPowerSetting enumeration from nl80211/nl80211.h:3856
const (
	TxPowerAutomatic = iota
	TxPowerLimited   = 1
	TxPowerFixed     = 2
)

// nl80211PacketPatternAttr as declared in nl80211/nl80211.h:3883
type nl80211PacketPatternAttr int32

// nl80211PacketPatternAttr enumeration from nl80211/nl80211.h:3883
const (
	__PktpatInvalid = iota
	PktpatMask      = 1
	PktpatPattern   = 2
	PktpatOffset    = 3
	Num_Pktpat      = 4
	Max_Pktpat      = Num_Pktpat - 1
)

// nl80211WowlanTriggers as declared in nl80211/nl80211.h:4011
type nl80211WowlanTriggers int32

// nl80211WowlanTriggers enumeration from nl80211/nl80211.h:4011
const (
	__WowlanTrigInvalid             = iota
	WowlanTrigAny                   = 1
	WowlanTrigDisconnect            = 2
	WowlanTrigMagicPkt              = 3
	WowlanTrigPktPattern            = 4
	WowlanTrigGtkRekeySupported     = 5
	WowlanTrigGtkRekeyFailure       = 6
	WowlanTrigEapIdentRequest       = 7
	WowlanTrig4wayHandshake         = 8
	WowlanTrigRfkillRelease         = 9
	WowlanTrigWakeupPkt80211        = 10
	WowlanTrigWakeupPkt80211Len     = 11
	WowlanTrigWakeupPkt8023         = 12
	WowlanTrigWakeupPkt8023Len      = 13
	WowlanTrigTcpConnection         = 14
	WowlanTrigWakeupTcpMatch        = 15
	WowlanTrigWakeupTcpConnlost     = 16
	WowlanTrigWakeupTcpNomoretokens = 17
	WowlanTrigNetDetect             = 18
	WowlanTrigNetDetectResults      = 19
	Num_WowlanTrig                  = 20
	Max_WowlanTrig                  = Num_WowlanTrig - 1
)

// nl80211WowlanTcpAttrs as declared in nl80211/nl80211.h:4129
type nl80211WowlanTcpAttrs int32

// nl80211WowlanTcpAttrs enumeration from nl80211/nl80211.h:4129
const (
	__WowlanTcpInvalid        = iota
	WowlanTcpSrcIpv4          = 1
	WowlanTcpDstIpv4          = 2
	WowlanTcpDstMac           = 3
	WowlanTcpSrcPort          = 4
	WowlanTcpDstPort          = 5
	WowlanTcpDataPayload      = 6
	WowlanTcpDataPayloadSeq   = 7
	WowlanTcpDataPayloadToken = 8
	WowlanTcpDataInterval     = 9
	WowlanTcpWakePayload      = 10
	WowlanTcpWakeMask         = 11
	Num_WowlanTcp             = 12
	Max_WowlanTcp             = Num_WowlanTcp - 1
)

// nl80211AttrCoalesceRule as declared in nl80211/nl80211.h:4174
type nl80211AttrCoalesceRule int32

// nl80211AttrCoalesceRule enumeration from nl80211/nl80211.h:4174
const (
	__CoalesceRuleInvalid      = iota
	AttrCoalesceRuleDelay      = 1
	AttrCoalesceRuleCondition  = 2
	AttrCoalesceRulePktPattern = 3
	Num_AttrCoalesceRule       = 4
	AttrCoalesceRuleMax        = Num_AttrCoalesceRule - 1
)

// nl80211CoalesceCondition as declared in nl80211/nl80211.h:4192
type nl80211CoalesceCondition int32

// nl80211CoalesceCondition enumeration from nl80211/nl80211.h:4192
const (
	CoalesceConditionMatch   = iota
	CoalesceConditionNoMatch = 1
)

// nl80211IfaceLimitAttrs as declared in nl80211/nl80211.h:4207
type nl80211IfaceLimitAttrs int32

// nl80211IfaceLimitAttrs enumeration from nl80211/nl80211.h:4207
const (
	IfaceLimitUnspec = iota
	IfaceLimitMax    = 1
	IfaceLimitTypes  = 2
	Num_IfaceLimit   = 3
	Max_IfaceLimit   = Num_IfaceLimit - 1
)

// nl80211IfCombinationAttrs as declared in nl80211/nl80211.h:4263
type nl80211IfCombinationAttrs int32

// nl80211IfCombinationAttrs enumeration from nl80211/nl80211.h:4263
const (
	IfaceCombUnspec             = iota
	IfaceCombLimits             = 1
	IfaceCombMaxnum             = 2
	IfaceCombStaApBiMatch       = 3
	IfaceCombNumChannels        = 4
	IfaceCombRadarDetectWidths  = 5
	IfaceCombRadarDetectRegions = 6
	Num_IfaceComb               = 7
	Max_IfaceComb               = Num_IfaceComb - 1
)

// nl80211PlinkState as declared in nl80211/nl80211.h:4296
type nl80211PlinkState int32

// nl80211PlinkState enumeration from nl80211/nl80211.h:4296
const (
	PlinkListen     = iota
	PlinkOpnSnt     = 1
	PlinkOpnRcvd    = 2
	PlinkCnfRcvd    = 3
	PlinkEstab      = 4
	PlinkHolding    = 5
	PlinkBlocked    = 6
	Num_PlinkStates = 7
	Max_PlinkStates = Num_PlinkStates - 1
)

// plinkActions as declared in nl80211/nl80211.h:4318
type plinkActions int32

// plinkActions enumeration from nl80211/nl80211.h:4318
const (
	PlinkActionNoAction = iota
	PlinkActionOpen     = 1
	PlinkActionBlock    = 2
	Num_PlinkActions    = 3
)

// nl80211RekeyData as declared in nl80211/nl80211.h:4340
type nl80211RekeyData int32

// nl80211RekeyData enumeration from nl80211/nl80211.h:4340
const (
	__RekeyDataInvalid = iota
	RekeyDataKek       = 1
	RekeyDataKck       = 2
	RekeyDataReplayCtr = 3
	Num_RekeyData      = 4
	Max_RekeyData      = Num_RekeyData - 1
)

// nl80211HiddenSsid as declared in nl80211/nl80211.h:4360
type nl80211HiddenSsid int32

// nl80211HiddenSsid enumeration from nl80211/nl80211.h:4360
const (
	HiddenSsidNotInUse     = iota
	HiddenSsidZeroLen      = 1
	HiddenSsidZeroContents = 2
)

// nl80211StaWmeAttr as declared in nl80211/nl80211.h:4376
type nl80211StaWmeAttr int32

// nl80211StaWmeAttr enumeration from nl80211/nl80211.h:4376
const (
	__StaWmeInvalid   = iota
	StaWmeUapsdQueues = 1
	StaWmeMaxSp       = 2
	__StaWmeAfterLast = 3
	StaWmeMax         = __StaWmeAfterLast - 1
)

// nl80211PmksaCandidateAttr as declared in nl80211/nl80211.h:4398
type nl80211PmksaCandidateAttr int32

// nl80211PmksaCandidateAttr enumeration from nl80211/nl80211.h:4398
const (
	__PmksaCandidateInvalid = iota
	PmksaCandidateIndex     = 1
	PmksaCandidateBssid     = 2
	PmksaCandidatePreauth   = 3
	Num_PmksaCandidate      = 4
	Max_PmksaCandidate      = Num_PmksaCandidate - 1
)

// nl80211TdlsOperation as declared in nl80211/nl80211.h:4417
type nl80211TdlsOperation int32

// nl80211TdlsOperation enumeration from nl80211/nl80211.h:4417
const (
	TdlsDiscoveryReq = iota
	TdlsSetup        = 1
	TdlsTeardown     = 2
	TdlsEnableLink   = 3
	TdlsDisableLink  = 4
)

// nl80211FeatureFlags as declared in nl80211/nl80211.h:4526
type nl80211FeatureFlags int32

// nl80211FeatureFlags enumeration from nl80211/nl80211.h:4526
const (
	FeatureSkTxStatus             = 1 << 0
	FeatureHtIbss                 = 1 << 1
	FeatureInactivityTimer        = 1 << 2
	FeatureCellBaseRegHints       = 1 << 3
	FeatureP2pDeviceNeedsChannel  = 1 << 4
	FeatureSae                    = 1 << 5
	FeatureLowPriorityScan        = 1 << 6
	FeatureScanFlush              = 1 << 7
	FeatureApScan                 = 1 << 8
	FeatureVifTxpower             = 1 << 9
	FeatureNeedObssScan           = 1 << 10
	FeatureP2pGoCtwin             = 1 << 11
	FeatureP2pGoOppps             = 1 << 12
	FeatureAdvertiseChanLimits    = 1 << 14
	FeatureFullApClientState      = 1 << 15
	FeatureUserspaceMpm           = 1 << 16
	FeatureActiveMonitor          = 1 << 17
	FeatureApModeChanWidthChange  = 1 << 18
	FeatureDsParamSetIeInProbes   = 1 << 19
	FeatureWfaTpcIeInProbes       = 1 << 20
	FeatureQuiet                  = 1 << 21
	FeatureTxPowerInsertion       = 1 << 22
	FeatureAcktoEstimation        = 1 << 23
	FeatureStaticSmps             = 1 << 24
	FeatureDynamicSmps            = 1 << 25
	FeatureSupportsWmmAdmission   = 1 << 26
	FeatureMacOnCreate            = 1 << 27
	FeatureTdlsChannelSwitch      = 1 << 28
	FeatureScanRandomMacAddr      = 1 << 29
	FeatureSchedScanRandomMacAddr = 1 << 30
	FeatureNdRandomMacAddr        = 1 << 31
)

// nl80211ExtFeatureIndex as declared in nl80211/nl80211.h:4595
type nl80211ExtFeatureIndex int32

// nl80211ExtFeatureIndex enumeration from nl80211/nl80211.h:4595
const (
	ExtFeatureVhtIbss          = iota
	ExtFeatureRrm              = 1
	ExtFeatureMuMimoAirSniffer = 2
	ExtFeatureScanStartTime    = 3
	ExtFeatureBssParentTsf     = 4
	ExtFeatureSetScanDwell     = 5
	ExtFeatureBeaconRateLegacy = 6
	ExtFeatureBeaconRateHt     = 7
	ExtFeatureBeaconRateVht    = 8
	Num_ExtFeatures            = 9
	Max_ExtFeatures            = Num_ExtFeatures - 1
)

// nl80211ProbeRespOffloadSupportAttr as declared in nl80211/nl80211.h:4625
type nl80211ProbeRespOffloadSupportAttr int32

// nl80211ProbeRespOffloadSupportAttr enumeration from nl80211/nl80211.h:4625
const (
	ProbeRespOffloadSupportWps    = 1 << 0
	ProbeRespOffloadSupportWps2   = 1 << 1
	ProbeRespOffloadSupportP2p    = 1 << 2
	ProbeRespOffloadSupport80211u = 1 << 3
)

// nl80211ConnectFailedReason as declared in nl80211/nl80211.h:4638
type nl80211ConnectFailedReason int32

// nl80211ConnectFailedReason enumeration from nl80211/nl80211.h:4638
const (
	ConnFailMaxClients    = iota
	ConnFailBlockedClient = 1
)

// nl80211ScanFlags as declared in nl80211/nl80211.h:4667
type nl80211ScanFlags int32

// nl80211ScanFlags enumeration from nl80211/nl80211.h:4667
const (
	ScanFlagLowPriority = 1 << 0
	ScanFlagFlush       = 1 << 1
	ScanFlagAp          = 1 << 2
	ScanFlagRandomAddr  = 1 << 3
)

// nl80211AclPolicy as declared in nl80211/nl80211.h:4687
type nl80211AclPolicy int32

// nl80211AclPolicy enumeration from nl80211/nl80211.h:4687
const (
	AclPolicyAcceptUnlessListed = iota
	AclPolicyDenyUnlessListed   = 1
)

// nl80211SmpsMode as declared in nl80211/nl80211.h:4702
type nl80211SmpsMode int32

// nl80211SmpsMode enumeration from nl80211/nl80211.h:4702
const (
	SmpsOff         = iota
	SmpsStatic      = 1
	SmpsDynamic     = 2
	__SmpsAfterLast = 3
	SmpsMax         = __SmpsAfterLast - 1
)

// nl80211RadarEvent as declared in nl80211/nl80211.h:4726
type nl80211RadarEvent int32

// nl80211RadarEvent enumeration from nl80211/nl80211.h:4726
const (
	RadarDetected    = iota
	RadarCacFinished = 1
	RadarCacAborted  = 2
	RadarNopFinished = 3
)

// nl80211DfsState as declared in nl80211/nl80211.h:4744
type nl80211DfsState int32

// nl80211DfsState enumeration from nl80211/nl80211.h:4744
const (
	DfsUsable      = iota
	DfsUnavailable = 1
	DfsAvailable   = 2
)

// nl80211ProtocolFeatures as declared in nl80211/nl80211.h:4758
type nl80211ProtocolFeatures int32

// nl80211ProtocolFeatures enumeration from nl80211/nl80211.h:4758
const (
	ProtocolFeatureSplitWiphyDump = 1 << 0
)

// nl80211CritProtoId as declared in nl80211/nl80211.h:4771
type nl80211CritProtoId int32

// nl80211CritProtoId enumeration from nl80211/nl80211.h:4771
const (
	CritProtoUnspec = iota
	CritProtoDhcp   = 1
	CritProtoEapol  = 2
	CritProtoApipa  = 3
	Num_CritProto   = 4
)

// nl80211RxmgmtFlags as declared in nl80211/nl80211.h:4790
type nl80211RxmgmtFlags int32

// nl80211RxmgmtFlags enumeration from nl80211/nl80211.h:4790
const (
	RxmgmtFlagAnswered = 1 << 0
)

// nl80211TdlsPeerCapability as declared in nl80211/nl80211.h:4824
type nl80211TdlsPeerCapability int32

// nl80211TdlsPeerCapability enumeration from nl80211/nl80211.h:4824
const (
	TdlsPeerHt  = 1 << 0
	TdlsPeerVht = 1 << 1
	TdlsPeerWmm = 1 << 2
)

// nl80211SchedScanPlan as declared in nl80211/nl80211.h:4843
type nl80211SchedScanPlan int32

// nl80211SchedScanPlan enumeration from nl80211/nl80211.h:4843
const (
	__SchedScanPlanInvalid   = iota
	SchedScanPlanInterval    = 1
	SchedScanPlanIterations  = 2
	__SchedScanPlanAfterLast = 3
	SchedScanPlanMax         = __SchedScanPlanAfterLast - 1
)

// nl80211BssSelectAttr as declared in nl80211/nl80211.h:4887
type nl80211BssSelectAttr int32

// nl80211BssSelectAttr enumeration from nl80211/nl80211.h:4887
const (
	__BssSelectAttrInvalid   = iota
	BssSelectAttrRssi        = 1
	BssSelectAttrBandPref    = 2
	BssSelectAttrRssiAdjust  = 3
	__BssSelectAttrAfterLast = 4
	BssSelectAttrMax         = __BssSelectAttrAfterLast - 1
)

// nl80211NanDualBandConf as declared in nl80211/nl80211.h:4907
type nl80211NanDualBandConf int32

// nl80211NanDualBandConf enumeration from nl80211/nl80211.h:4907
const (
	NanBandDefault = 1 << 0
	NanBand2ghz    = 1 << 1
	NanBand5ghz    = 1 << 2
)

// nl80211NanFunctionType as declared in nl80211/nl80211.h:4922
type nl80211NanFunctionType int32

// nl80211NanFunctionType enumeration from nl80211/nl80211.h:4922
const (
	NanFuncPublish         = iota
	NanFuncSubscribe       = 1
	NanFuncFollowUp        = 2
	__NanFuncTypeAfterLast = 3
	NanFuncMaxType         = __NanFuncTypeAfterLast - 1
)

// nl80211NanPublishType as declared in nl80211/nl80211.h:4940
type nl80211NanPublishType int32

// nl80211NanPublishType enumeration from nl80211/nl80211.h:4940
const (
	NanSolicitedPublish   = 1 << 0
	NanUnsolicitedPublish = 1 << 1
)

// nl80211NanFuncTermReason as declared in nl80211/nl80211.h:4954
type nl80211NanFuncTermReason int32

// nl80211NanFuncTermReason enumeration from nl80211/nl80211.h:4954
const (
	NanFuncTermReasonUserRequest = iota
	NanFuncTermReasonTtlExpired  = 1
	NanFuncTermReasonError       = 2
)

// nl80211NanFuncAttributes as declared in nl80211/nl80211.h:5006
type nl80211NanFuncAttributes int32

// nl80211NanFuncAttributes enumeration from nl80211/nl80211.h:5006
const (
	__NanFuncInvalid       = iota
	NanFuncType            = 1
	NanFuncServiceId       = 2
	NanFuncPublishType     = 3
	NanFuncPublishBcast    = 4
	NanFuncSubscribeActive = 5
	NanFuncFollowUpId      = 6
	NanFuncFollowUpReqId   = 7
	NanFuncFollowUpDest    = 8
	NanFuncCloseRange      = 9
	NanFuncTtl             = 10
	NanFuncServiceInfo     = 11
	NanFuncSrf             = 12
	NanFuncRxMatchFilter   = 13
	NanFuncTxMatchFilter   = 14
	NanFuncInstanceId      = 15
	NanFuncTermReason      = 16
	Num_NanFuncAttr        = 17
	NanFuncAttrMax         = Num_NanFuncAttr - 1
)

// nl80211NanSrfAttributes as declared in nl80211/nl80211.h:5045
type nl80211NanSrfAttributes int32

// nl80211NanSrfAttributes enumeration from nl80211/nl80211.h:5045
const (
	__NanSrfInvalid = iota
	NanSrfInclude   = 1
	NanSrfBf        = 2
	NanSrfBfIdx     = 3
	NanSrfMacAddrs  = 4
	Num_NanSrfAttr  = 5
	NanSrfAttrMax   = Num_NanSrfAttr - 1
)

// nl80211NanMatchAttributes as declared in nl80211/nl80211.h:5070
type nl80211NanMatchAttributes int32

// nl80211NanMatchAttributes enumeration from nl80211/nl80211.h:5070
const (
	__NanMatchInvalid = iota
	NanMatchFuncLocal = 1
	NanMatchFuncPeer  = 2
	Num_NanMatchAttr  = 3
	NanMatchAttrMax   = Num_NanMatchAttr - 1
)
