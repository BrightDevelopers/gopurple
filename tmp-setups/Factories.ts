import { deepCopy } from "@sys/Tools";
import { SetupPackage } from "./SetupPackage";
import { UserCredentials } from "@api/Credentials";
import SetupToken from "@v202206/DTO/SetupToken";

const defaultBsnCloudSetupPackageV2 = require('./DefaultBsnCloudSetupPackageV2.json');
const defaultBsnCloudSetupPackageV3 = require('./DefaultBsnCloudSetupPackageV3.json');
const defaultPartnerApplicationSetupPackageV3 = require('./DefaultPartnerApplicationSetupPackageV3.json');
const defaultSimpleFileNetworkingSetupPackageV3 = require('./DefaultSimpleFileNetworkingSetupPackageV3.json');
const defaultLocalFileNetworkingSetupPackageV3 = require('./DefaultLocalFileNetworkingSetupPackageV3.json');

export class SetupPackageFactory {
	public static bsnCloudPackage(user: UserCredentials, token: SetupToken, packageName: string): SetupPackage {
		const clone = deepCopy(defaultBsnCloudSetupPackageV3) as SetupPackage;
		clone.bDeploy.networkName = user.network!;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}

	public static bsnCloudPackageV2(user: UserCredentials, token: SetupToken, packageName: string): SetupPackage {
		const clone = deepCopy(defaultBsnCloudSetupPackageV2) as SetupPackage;
		clone.bDeploy.networkName = user.network!;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}

	public static partnerApplicationPackage(user: UserCredentials, token: SetupToken, packageName: string, url: string): SetupPackage {
		const clone = deepCopy(defaultPartnerApplicationSetupPackageV3) as SetupPackage;
		clone.bDeploy.networkName = user.network!;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bDeploy.url = url;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}

	public static partnerApplicationPackageWithoutNetwork(user: UserCredentials, token: SetupToken, packageName: string, url: string): SetupPackage {
		const clone = deepCopy(defaultPartnerApplicationSetupPackageV3) as SetupPackage;
		// networkName intentionally omitted for third-party CMS scenarios
		clone.bDeploy.networkName = undefined as any;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bDeploy.url = url;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}

	public static simpleFileNetworkingPackage(user: UserCredentials, token: SetupToken, packageName: string): SetupPackage {
		const clone = deepCopy(defaultSimpleFileNetworkingSetupPackageV3) as SetupPackage;
		clone.bDeploy.networkName = user.network!;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}

	public static localFileNetworkingPackage(user: UserCredentials, token: SetupToken, packageName: string): SetupPackage {
		const clone = deepCopy(defaultLocalFileNetworkingSetupPackageV3) as SetupPackage;
		clone.bDeploy.networkName = user.network!;
		clone.bDeploy.packageName = packageName;
		clone.bDeploy.username = user.username;
		clone.bsnDeviceRegistrationTokenEntity = {
			token: token.token,
			validFrom: token.validFrom,
			validTo: token.validTo,
			scope: token.scope
		};
		return clone;
	}
}
