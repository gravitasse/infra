#include <stdio.h>
#include "onlp.h"
#include <onlp/fan.h>
#include <onlp/sys.h>

static onlp_oid_t fan_oid_table[ONLP_OID_TABLE_SIZE];
static onlp_fan_info_t fan_info_table[ONLP_OID_TABLE_SIZE];
static int flag[ONLP_OID_TABLE_SIZE] = {0};
onlp_sys_info_t si;

int Init() {
	int ret = 0;
	int i = 0;
	onlp_oid_t* oidp;

	printf("=============Init================");
	onlp_init();
	ret = onlp_sys_info_get(&si);
	if (ret < 0) {
		printf("onlp_sys_info_get() failed during init. Return Value:%d\n", ret);
		return ret;
	}

	ONLP_OID_TABLE_ITER_TYPE(si.hdr.coids, oidp, FAN) {
		fan_oid_table[i++] = *oidp;
	}
	
	return 0;
}

int GetMaxNumOfFans() {
	return ONLP_OID_TABLE_SIZE;
}

int GetAllFanState(fan_info_t *info, int count) {
	printf("Count = %d", count);
	int i = 0;
	for (i = 0; i < count; i++) {
		onlp_fan_info_t fi;
		int fid = ONLP_OID_ID_GET(fan_oid_table[i]);
		info[i].FanId = fid;
		info[i].valid = 0;
		if (fan_oid_table[i] == 0) {
			continue;
		}
		if (onlp_fan_info_get(fan_oid_table[i], &fi) < 0) {
			printf("Failure retreiving status of Fan Id %d\n", fid);
			continue;
		}
		printf("Fan Id : %d Fan Status : %u\n", fid, fi.status);
		info[i].valid = 1;
		if (!(fi.status & 0x1)) {
			info[i].Status = FAN_STATUS_MISSING;
		} else if (fi.status & ONLP_FAN_STATUS_FAILED) {
			info[i].Status = FAN_STATUS_FAILED;
		} else if (fi.status & 0x1) {
			info[i].Status = FAN_STATUS_PRESENT;
		}

		info[i].Speed = fi.rpm;
		info[i].Direction = (fi.status & (1 << 2)) ? FAN_DIR_B2F : FAN_DIR_F2B;
		if (fi.rpm > 0) {
			info[i].Mode = FAN_MODE_ON;
		} else {
			info[i].Mode = FAN_MODE_OFF;
		}

		strncpy(info[i].Model, fi.model, 100);
		strncpy(info[i].SerialNum, fi.serial, 100);
	}
	return 0;
}

int GetFanState(fan_info_t *info, int fanId) {
	int i = 0;
	for (i = 0; i < sizeof(fan_oid_table)/sizeof(fan_oid_table[0]); i++) {
		onlp_fan_info_t fi;
		int fid = ONLP_OID_ID_GET(fan_oid_table[i]);
		if (fid != fanId) {
			continue;
		}
		info[0].FanId = fid;
		info[0].valid = 0;
		if (fan_oid_table[i] == 0) {
			continue;
		}
		if (onlp_fan_info_get(fan_oid_table[i], &fi) < 0) {
			printf("Failure retreiving status of Fan Id %d\n", fid);
			continue;
		}
		printf("Fan Id : %d Fan Status : %u\n", fid, fi.status);
		info[0].valid = 1;
		if (!(fi.status & 0x1)) {
			info[0].Status = FAN_STATUS_MISSING;
		} else if (fi.status & ONLP_FAN_STATUS_FAILED) {
			info[0].Status = FAN_STATUS_FAILED;
		} else if (fi.status & 0x1) {
			info[0].Status = FAN_STATUS_PRESENT;
		}

		info[0].Speed = fi.rpm;
		info[0].Direction = (fi.status & (1 << 2)) ? FAN_DIR_B2F : FAN_DIR_F2B;
		if (fi.rpm > 0) {
			info[0].Mode = FAN_MODE_ON;
		} else {
			info[0].Mode = FAN_MODE_OFF;
		}

		strncpy(info[0].Model, fi.model, 100);
		strncpy(info[0].SerialNum, fi.serial, 100);
		return 0;
	}
	return -1;
}
