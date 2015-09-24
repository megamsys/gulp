/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package gulpd

import (

)


/*func (s *S) TestInsertEmptyContainerInDBForward(c *check.C) {
	app := provisiontest.NewFakeApp("myapp", "python", 1)
	args := runContainerActionsArgs{
		app:           app,
		imageID:       "image-id",
		buildingImage: "next-image",
		provisioner:   s.p,
	}
	context := action.FWContext{Params: []interface{}{args}}
	r, err := insertEmptyContainerInDB.Forward(context)
	c.Assert(err, check.IsNil)
	cont := r.(container.Container)
	c.Assert(cont, check.FitsTypeOf, container.Container{})
	c.Assert(cont.AppName, check.Equals, app.GetName())
	c.Assert(cont.Type, check.Equals, app.GetPlatform())
	c.Assert(cont.Name, check.Not(check.Equals), "")
	c.Assert(strings.HasPrefix(cont.Name, app.GetName()+"-"), check.Equals, true)
	c.Assert(cont.Name, check.HasLen, 26)
	c.Assert(cont.Status, check.Equals, "created")
	c.Assert(cont.Image, check.Equals, "image-id")
	c.Assert(cont.BuildingImage, check.Equals, "next-image")
	coll := s.p.Collection()
	defer coll.Close()
	defer coll.Remove(bson.M{"name": cont.Name})
	var retrieved container.Container
	err = coll.Find(bson.M{"name": cont.Name}).One(&retrieved)
	c.Assert(err, check.IsNil)
	c.Assert(retrieved.Name, check.Equals, cont.Name)
}*/


