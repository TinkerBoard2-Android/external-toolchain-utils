# Copyright 2017 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

LOCAL_PATH := $(call my-dir)

include $(CLEAR_VARS)
LOCAL_MODULE_TAGS := tests
LOCAL_C_INCLUDES := $(LOCAL_PATH)/source
LOCAL_SRC_FILES:= apps/synthmark.cpp
LOCAL_CFLAGS += -g -std=c++11 -Ofast
LOCAL_CFLAGS += $(CFLAGS_FOR_BENCH_SUITE)
LOCAL_LDFLAGS += $(LDFLAGS_FOR_BENCH_SUITE)
#LOCAL_SHARED_LIBRARIES := libcutils libutils
LOCAL_MODULE := synthmark
include $(BUILD_EXECUTABLE)
