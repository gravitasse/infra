//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//
#ifndef PLUGINCOMMON_H
#define PLUGINCOMMON_H

#include <stdio.h>

#define DEFAULT_SIZE 128

typedef enum fan_dir_e {
	FAN_DIR_B2F,
	FAN_DIR_F2B,
	FAN_DIR_INVALID,
} fan_dir_t;

typedef enum fan_mode_e {
    FAN_MODE_OFF,
    FAN_MODE_ON,
} fan_mode_t;

typedef enum fan_status_e {
    FAN_STATUS_PRESENT,
    FAN_STATUS_MISSING,
    FAN_STATUS_FAILED,
    FAN_STATUS_NORMAL,
} fan_status_t;


typedef struct fan_info {
	int valid;
	int FanId;
	fan_mode_t Mode;
	int Speed;
	fan_dir_t Direction;
	fan_status_t Status;
	char Model[DEFAULT_SIZE];
	char SerialNum[DEFAULT_SIZE];
} fan_info_t;

typedef struct sys_info {
    char product_name[DEFAULT_SIZE];
    char serial_number[DEFAULT_SIZE];
    char manufacturer[DEFAULT_SIZE];
    char vendor[DEFAULT_SIZE];
    char platform_name[DEFAULT_SIZE];
    char onie_version[DEFAULT_SIZE];
    char label_revision[DEFAULT_SIZE];
} sys_info_t;

typedef enum {
    SFP_ERROR = -1,
    SFP_OK = 0,
    SFP_MISSING = 1,
}SFP_RET;

typedef struct sfp_info {
   int sfp_id;
   unsigned int spf_speed; /* in Mbps */
   int sfp_present;
   int sfp_los;
   char serial_number[12];
   char eeprom[256];
} sfp_info_t;


#endif // PLUGINCOMMON_H
