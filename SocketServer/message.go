/* Copyright (C) 2013 Christopher Lillthors and Viktor Kronvall
 * This file is part of GoWebSocket.
 *
 * GoWebSocket is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GoWebSocket is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GoWebSocket.  If not, see <http://www.gnu.org/licenses/>.
 */

package SocketServer

type MessageStruct struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}
