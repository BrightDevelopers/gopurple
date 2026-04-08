import { DateTime } from "@sys/Types";

export interface SetupPackage {
	version: string;
	bDeploy: {
		networkName?: string,
		username: string,
		client: string,
		packageName: string,
		url?: string
	};
	firmwareUpdatesByFamily: Record<string, any>;
	firmwareUpdateType: string;
	setupType: string;
	bsnDeviceRegistrationTokenEntity: {
		token: string;
		scope: string;
		validFrom: DateTime;
		validTo: DateTime;
	};
	enableSerialDebugging: boolean;
	enableSystemLogDebugging: boolean;
	remoteDwsEnabled: boolean;
	dwsEnabled: boolean;
	dwsPassword: string;
	dwsPasswordPreviousSavedTimeStamp: number;
	lwsEnabled: boolean;
	lwsConfig: string;
	lwsUserName: string;
	lwsPassword: string;
	lwsEnableUpdateNotifications: boolean;
	bsnCloudEnabled: boolean;
	deviceName: string;
	deviceDescription: string;
	unitNamingMethod: string;
	timeZone: string;
	bsnGroupName: string;
	timeBetweenNetConnects: number;
	timeBetweenHeartbeats: number;
	sfnWebFolderUrl: string;
	sfnUserName: string;
	sfnPassword: string;
	sfnEnableBasicAuthentication: boolean;
	playbackLoggingEnabled: boolean;
	eventLoggingEnabled: boolean;
	diagnosticLoggingEnabled: boolean;
	stateLoggingEnabled: boolean;
	variableLoggingEnabled: boolean;
	uploadLogFilesAtBoot: boolean;
	uploadLogFilesAtSpecificTime: boolean;
	uploadLogFilesTime: number;
	logHandlerUrl: string;
	enableRemoteSnapshot: boolean;
	remoteSnapshotInterval: number;
	remoteSnapshotMaxImages: number;
	remoteSnapshotJpegQualityLevel: number;
	remoteSnapshotScreenOrientation: string;
	remoteSnapshotHandlerUrl: string;
	idleScreenColor: {
		r: number;
		g: number;
		b: number;
		a: number;
	};
	networkDiagnosticsEnabled: boolean;
	testEthernetEnabled: boolean;
	testWirelessEnabled: boolean;
	testInternetEnabled: boolean;
	useCustomSplashScreen: boolean;
	BrightWallName: string;
	BrightWallScreenNumber: string;
	contentDownloadsRestricted: boolean;
	contentDownloadRangeStart: number;
	contentDownloadRangeEnd: number;
	heartbeatsRestricted: boolean;
	heartbeatsRangeStart: number;
	heartbeatsRangeEnd: number;
	usbUpdatePassword: string;
	inheritNetworkProperties: boolean;
	internalCaArtifacts: any[];
	network: {
		timeServers: string[];
		hostname: null | string;
		dns: null | string;
		proxyServer: null | string;
		proxyBypass: null | string;
		interfaces: {
			id: string;
			name: string;
			metric: null | number;
			type: string;
			proto: string;
			ip: string[];
			gateway: null | string;
			dns: string[];
			showInUi: boolean;
			rateLimitDuringInitialDownloads: null | number;
			rateLimitInsideContentDownloadWindow: null | number;
			rateLimitOutsideContentDownloadWindow: null | number;
			contentDownloadEnabled: boolean;
			textFeedsDownloadEnabled: boolean;
			mediaFeedsDownloadEnabled: boolean;
			healthReportingEnabled: boolean;
			logsUploadEnabled: boolean;
			wpaSettings: {
				enableWPAEnterpriseAuthentication: boolean;
				wpaEnterpriseVariant: string;
				eapCertificateType: string;
				eapCertificateFile: null | string;
				eapCertificatePassphrase: string;
				eapPemOrDerKeyFile: null | string;
				peapUsername: string;
				peapPassphrase: string;
				caCertificateFile: null | string;
			};
		}[];
		certificateName: string;
		certificateFile: null | string;
	};
}
