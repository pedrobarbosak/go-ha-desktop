// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {config} from '../models';
import {ha} from '../models';

export function GetConfig():Promise<config.Config>;

export function GetDevices():Promise<Array<ha.Device>>;

export function GetError():Promise<void>;

export function SetBrightness(arg1:string,arg2:number):Promise<ha.Device>;

export function TestConnection(arg1:string,arg2:string):Promise<void>;

export function TurnOff(arg1:string):Promise<ha.Device>;

export function TurnOn(arg1:string):Promise<ha.Device>;

export function UpdateConfig(arg1:config.Config):Promise<void>;
