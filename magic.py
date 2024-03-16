import re
import os
import glob
import pathlib

template_model = """package model

// {{MODEL_NAME}} {{MODEL_NAME_HUMAN}}
type {{MODEL_NAME}} struct {
	Base
	DataOwner
	{{MODEL_NAME}}API
}

// {{MODEL_NAME}}API {{MODEL_NAME_HUMAN}} API
type {{MODEL_NAME}}API struct {
{{MODEL_PROPERTIES}}
}
"""

template_get = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
)

// Get{{MODEL_NAME}} godoc
// @Summary List of {{MODEL_NAME_HUMAN}}
// @Description List of {{MODEL_NAME_HUMAN}}
// @Param page query int false "Page number start from zero"
// @Param size query int false "Size per page, default `0`"
// @Param sort query string false "Sort by field, adding dash (`-`) at the beginning means descending and vice versa"
// @Param fields query string false "Select specific fields with comma separated"
// @Param filters query string false "custom filters, see [more details](https://github.com/morkid/paginate#filter-format)"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} lib.Page{items=[]model.{{MODEL_NAME}}} "List of {{MODEL_NAME_HUMAN}}"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /{{ROUTER_NAME}} [get]
// @Tags {{MODEL_NAME}}
func Get{{MODEL_NAME}}(c *fiber.Ctx) error {
    db := services.DB.WithContext(c.UserContext()).WithContext(c.UserContext())
	pg := services.PG

	mod := db.Model(&model.{{MODEL_NAME}}{})

	page := pg.With(mod).Request(c.Request()).Response(&[]model.{{MODEL_NAME}}{})

	return lib.OK(c, page)
}
"""

template_get_id = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Get{{MODEL_NAME}}ID godoc
// @Summary Get a {{MODEL_NAME_HUMAN}} by id
// @Description Get a {{MODEL_NAME_HUMAN}} by id
// @Param id path string true "{{MODEL_NAME_HUMAN}} ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.{{MODEL_NAME}} "{{MODEL_NAME_HUMAN}} data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /{{ROUTER_NAME}}/{id} [get]
// @Tags {{MODEL_NAME}}
func Get{{MODEL_NAME}}ID(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.{{MODEL_NAME}}
	result := db.Model(&data).
		Where(db.Where(model.{{MODEL_NAME}}{
			Base: model.Base{
				ID: &id,
			},
		})).
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	return lib.OK(c, data)
}
"""

template_delete = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
)

// Delete{{MODEL_NAME}} godoc
// @Summary Delete {{MODEL_NAME_HUMAN}} by id
// @Description Delete {{MODEL_NAME_HUMAN}} by id
// @Param id path string true "{{MODEL_NAME_HUMAN}} ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /{{ROUTER_NAME}}/{id} [delete]
// @Tags {{MODEL_NAME}}
func Delete{{MODEL_NAME}}(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())

	var data model.{{MODEL_NAME}}
	result := db.Model(&data).Where("id = ?", c.Params("id")).Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	db.Delete(&data)

	return lib.OK(c)
}
"""

template_post = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
)

// Post{{MODEL_NAME}} godoc
// @Summary Create new {{MODEL_NAME_HUMAN}}
// @Description Create new {{MODEL_NAME_HUMAN}}
// @Param data body model.{{MODEL_NAME}}API true "{{MODEL_NAME_HUMAN}} data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.{{MODEL_NAME}} "{{MODEL_NAME_HUMAN}} data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /{{ROUTER_NAME}} [post]
// @Tags {{MODEL_NAME}}
func Post{{MODEL_NAME}}(c *fiber.Ctx) error {
	api := new(model.{{MODEL_NAME}}API)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	db := services.DB.WithContext(c.UserContext())

	var data model.{{MODEL_NAME}}
	lib.Merge(api, &data)
	data.CreatorID = lib.GetXUserID(c)

	if err := db.Create(&data).Error; nil != err {
		return lib.ErrorConflict(c, err.Error())
	}

	return lib.Created(c, data)
}
"""

template_put = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Put{{MODEL_NAME}} godoc
// @Summary Update {{MODEL_NAME_HUMAN}} by id
// @Description Update {{MODEL_NAME_HUMAN}} by id
// @Param id path string true "{{MODEL_NAME_HUMAN}} ID"
// @Param data body model.{{MODEL_NAME}}API true "{{MODEL_NAME_HUMAN}} data"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.{{MODEL_NAME}} "{{MODEL_NAME_HUMAN}} data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /{{ROUTER_NAME}}/{id} [put]
// @Tags {{MODEL_NAME}}
func Put{{MODEL_NAME}}(c *fiber.Ctx) error {
	api := new(model.{{MODEL_NAME}}API)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.{{MODEL_NAME}}
	result := db.Model(&data).
		Where(db.Where(model.{{MODEL_NAME}}{
			Base: model.Base{
				ID: &id,
			},
		})).
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	lib.Merge(api, &data)
	data.ModifierID = lib.GetXUserID(c)

	if err := db.Model(&data).Updates(&data).Error; nil != err {
		return lib.ErrorConflict(c, err.Error())
	}

	return lib.OK(c, data)
}
"""

template_get_test = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestGet{{MODEL_NAME}}(t *testing.T) {
	db := services.DBConnectTest()
	lib.LoadEnvironment(config.Environment)

	app := fiber.New()
	app.Use(middleware.TokenValidator())

	app.Get("/{{ROUTER_NAME}}", Get{{MODEL_NAME}})

	initial := model.{{MODEL_NAME}}{
		{{MODEL_NAME}}API: model.{{MODEL_NAME}}API{
			{{MODEL_VALUES}}
		},
	}

	db.Create(&initial)

	headers := map[string]string{
		viper.GetString("HEADER_TOKEN_KEY"): viper.GetString("VALUE_TOKEN_KEY"),
	}

	uri := "/{{ROUTER_NAME}}?page=0&size=1"
	response, body, err := lib.GetTest(app, uri, headers)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 200, response.StatusCode, "getting response code")
	utils.AssertEqual(t, false, nil == body, "validate response body")
	utils.AssertEqual(t, float64(1), body["total"], "getting response body")


    // test invalid token
	response, _, err = lib.GetTest(app, uri, nil)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 401, response.StatusCode, "getting response code")


	sqlDB, _ := db.DB()
	sqlDB.Close()
}
"""

template_get_id_test = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestGet{{MODEL_NAME}}ID(t *testing.T) {
	db := services.DBConnectTest()
	lib.LoadEnvironment(config.Environment)
    
	app := fiber.New()
	app.Use(middleware.TokenValidator())
    
	app.Get("/{{ROUTER_NAME}}/:id", Get{{MODEL_NAME}}ID)

	initial := model.{{MODEL_NAME}}{
		{{MODEL_NAME}}API: model.{{MODEL_NAME}}API{
			{{MODEL_VALUES}}
		},
	}

	db.Create(&initial)

	headers := map[string]string{
		viper.GetString("HEADER_TOKEN_KEY"): viper.GetString("VALUE_TOKEN_KEY"),
	}

	uri := "/{{ROUTER_NAME}}/" + initial.ID.String()
	response, body, err := lib.GetTest(app, uri, headers)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 200, response.StatusCode, "getting response code")
	utils.AssertEqual(t, false, nil == body, "validate response body")
	utils.AssertEqual(t, initial.ID.String(), body["id"], "getting response body")

	// test get non existing id
	uri = "/{{ROUTER_NAME}}/non-existing-id"
	response, _, err = lib.GetTest(app, uri, headers)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 404, response.StatusCode, "getting response code")

	// test invalid token
	response, _, err = lib.GetTest(app, uri, nil)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 401, response.StatusCode, "getting response code")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
"""

template_post_test = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
    "github.com/spf13/viper"
)

func TestPost{{MODEL_NAME}}(t *testing.T) {
	db := services.DBConnectTest()
	lib.LoadEnvironment(config.Environment)
    
	app := fiber.New()
	app.Use(middleware.TokenValidator())

	app.Post("/{{ROUTER_NAME}}", Post{{MODEL_NAME}})

	uri := "/{{ROUTER_NAME}}"

	payload := `{
		{{JSON_VALUES}}
	}`

	headers := map[string]string{
		"Content-Type": "application/json",
		viper.GetString("HEADER_TOKEN_KEY"): viper.GetString("VALUE_TOKEN_KEY"),
	}

	response, body, err := lib.PostTest(app, uri, headers, payload)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 201, response.StatusCode, "getting response code")
	utils.AssertEqual(t, false, nil == body, "validate response body")

	// test invalid json format
	response, _, err = lib.PostTest(app, uri, headers, "invalid json format")
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 400, response.StatusCode, "getting response code")

    // test duplicate data
    response, _, err = lib.PostTest(app, uri, headers, payload)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 409, response.StatusCode, "getting response code")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
"""

template_put_test = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
    "github.com/spf13/viper"
)

func TestPut{{MODEL_NAME}}(t *testing.T) {
	db := services.DBConnectTest()
	lib.LoadEnvironment(config.Environment)
    
	app := fiber.New()
	app.Use(middleware.TokenValidator())

	app.Put("/{{ROUTER_NAME}}/:id", Put{{MODEL_NAME}})

	initial := model.{{MODEL_NAME}}{
		{{MODEL_NAME}}API: model.{{MODEL_NAME}}API{
			{{MODEL_VALUES}}
		},
	}

	initial2 := model.{{MODEL_NAME}}{
		{{MODEL_NAME}}API: model.{{MODEL_NAME}}API{
			{{MODEL_VALUES}}
		},
	}

	db.Create(&initial)
	db.Create(&initial2)

	uri := "/{{ROUTER_NAME}}/" + initial.ID.String()

	payload := `{
		{{JSON_VALUES}}
	}`

	headers := map[string]string{
		"Content-Type": "application/json",
		viper.GetString("HEADER_TOKEN_KEY"): viper.GetString("VALUE_TOKEN_KEY"),
	}

	response, body, err := lib.PutTest(app, uri, headers, payload)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 200, response.StatusCode, "getting response code")
	utils.AssertEqual(t, false, nil == body, "validate response body")

	// test invalid json body
	response, _, err = lib.PutTest(app, uri, headers, "invalid json format")
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 400, response.StatusCode, "getting response code")

	// test update with non existing id
	uri = "/{{ROUTER_NAME}}/non-existing-id"
	response, _, err = lib.PutTest(app, uri, headers, payload)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 404, response.StatusCode, "getting response code")

    // test duplicate data
    uri = "/{{ROUTER_NAME}}/" + initial2.ID.String()
	response, _, err = lib.PutTest(app, uri, headers, payload)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 409, response.StatusCode, "getting response code")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
"""

template_delete_test = """package {{PACKAGE_NAME}}

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestDelete{{MODEL_NAME}}(t *testing.T) {
	db := services.DBConnectTest()
	lib.LoadEnvironment(config.Environment)
    
	app := fiber.New()
	app.Use(middleware.TokenValidator())

	app.Delete("/{{ROUTER_NAME}}/:id", Delete{{MODEL_NAME}})

	initial := model.{{MODEL_NAME}}{
		{{MODEL_NAME}}API: model.{{MODEL_NAME}}API{
			{{MODEL_VALUES}}
		},
	}

	db.Create(&initial)

	headers := map[string]string{
		viper.GetString("HEADER_TOKEN_KEY"): viper.GetString("VALUE_TOKEN_KEY"),
	}

	uri := "/{{ROUTER_NAME}}/" + initial.ID.String()
	response, _, err := lib.DeleteTest(app, uri, headers)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 200, response.StatusCode, "getting response code")

	// test delete with non existing id
	response, _, err = lib.DeleteTest(app, uri, headers)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 404, response.StatusCode, "getting response code")


	// test invalid token
	response, _, err = lib.DeleteTest(app, uri, nil)
	utils.AssertEqual(t, nil, err, "sending request")
	utils.AssertEqual(t, 401, response.StatusCode, "getting response code")

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
"""

def create_dir(dir_name):
    pathlib.Path(dir_name).mkdir(parents=True, exist_ok=True)

def pascal_case(value):
    new_value = ""
    value = re.sub(r"[^a-zA-Z0-9]+", "_", value)
    for val in value.split("_"):
        if val == "":
            continue
        next_val = ""
        if len(val) > 1:
            next_val = val[1:]

        if val == "id":
            next_val = "D"
        new_value = new_value + val[0].upper() + next_val
    return new_value

def capitalize(value):
    values = []
    value = re.sub(r"[^a-zA-Z0-9]+", "_", value)
    for val in value.split("_"):
        if val == "":
            continue
        next_val = ""
        if len(val) > 1:
            next_val = val[1:]

        if val == "id":
            next_val = "D"
        values.append(val[0].upper() + next_val)
    return " ".join(values)

def snake_case(value):
    values = []
    value = re.sub(r"[^a-zA-Z0-9]+", "_", value)
    for val in value.split("_"):
        values.append(val.lower())

    return "_".join(values)

def add_fields(size, tables):
    print("------------------------------")
    extra = ""
    if size > 0:
        extra = " [press enter to finish]"
    field_name = input("Create new field name"+extra+": ")
    field_type = ""
    comment = ""
    enable_unique = False
    enable_null = True
    skip_example = False
    tags = []
    gorms = []
    if field_name != "":
        default_field_type = ""
        extra_type = ""
        if re.match(r".*_id$", field_name):
            default_field_type = "uuid"
            extra_type = " (default uuid)"
        elif re.match(r"^(is_.*|.+able)$", field_name):
            default_field_type = "bool"
            extra_type = " (default bool)"
        field_type = input("Field type [string|int|int64|float|float64|bool|uuid|date|date-time|etc]"+extra_type+": ")
        if field_type == "" and default_field_type != "":
            field_type = default_field_type
        tags.append('json:"'+snake_case(field_name)+',omitempty"')
        if field_type == "":
            field_type = "string"
        if field_type == "uuid":
            tags.append('swaggertype:"string" format:"uuid"')
            field_type = "uuid.UUID"
            skip_example = True
        elif field_type == "date":
            tags.append('format:"date" swaggertype:"string"')
            gorms.append('type:date')
            field_type = "strfmt.Date"
            skip_example = True
        elif field_type == "int":
            gorms.append('type:smallint')
            skip_example = True
        elif field_type == "date-time":
            tags.append('format:"date-time" swaggertype:"string"')
            gorms.append('type:timestamptz')
            field_type = "strfmt.DateTime"
            skip_example = True
        elif field_type == "string":
            print("> Notes: set length greater than or equal 4000 to set sql field as text")
            field_length = input("Field length (default 256): ")
            if field_length != "":
                if int(field_length) >= 4000:
                    gorms.append('type:text')
                else:
                    gorms.append('type:varchar('+field_length+')')
            else:
                gorms.append('type:varchar(256)')
        if not skip_example:
            example = input("Example value: ")
            if example != "":
                tags.append('example:"'+example+'"')
        default_comment = capitalize(field_name)
        field_comment = input("Description (default '"+default_comment+"'): ")
        if field_comment != "":
            comment = " // " + field_comment
        else:
            comment = " // " + default_comment

        unique_field = input("Set unique [y/n] (default n): ")
        enable_unique = unique_field != "" and unique_field.lower() == "y"
        if enable_unique:
            unique_name = input("Unique Name (default `idx_"+ tables +"_"+ field_name +"_unique) :`")
            if unique_name is None or unique_name == "":
                unique_name = "idx_"+ tables + "_" + field_name + "_unique"
            gorms.append("index:"+unique_name+",unique,where:deleted_at is null;not null")
        else:
            null_field = input("Allow null [y/n] (default y): ")
            enable_null = null_field != "" and null_field.lower() == "n"
            if enable_null:
                gorms.append('not null')
    
    if len(gorms) > 0:
        tags.append('gorm:"'+(";".join(gorms))+'"')

    tag = " ".join(tags)
    return {
        "lower_name": snake_case(field_name),
        "name": pascal_case(field_name),
        "type": "*"+field_type,
        "tag": "`" + tag + "`",
        "comment": comment,
        "unique": enable_unique,
    }

def kebab_case(value):
    values = []
    value = re.sub(r"[^a-zA-Z0-9]+", "-", value)
    for val in value.split("-"):
        if val == "":
            continue
        values.append(val.lower())

    return "-".join(values)

def lower_alphanumeric(value):
    return re.sub(r"[^a-zA-Z0-9]+", "", value).lower()

def clean_code(file_name):
    # os.system("goimports -w "+file_name)
    os.system("gofmt -w "+file_name)

def create_endpoints():
    module_name = input("Module name: ")
    default_router_name = kebab_case(module_name)
    if len(default_router_name) > 0:
        if default_router_name[-1] == "s":
            default_router_name = default_router_name + "es"
        elif default_router_name[-1] == "y":
            default_router_name = default_router_name[0:-1] + "ies"
        else:
            default_router_name = default_router_name + "s"
    extra = ""
    if default_router_name != "":
        extra = " (/"+default_router_name+")"
    router_name = input("router name"+extra+": ")
    if router_name == "":
        router_name = default_router_name
    router_name = re.sub(r"^/+", "", router_name)
    router_name = re.sub(r"/+$", "", router_name)
    router_name = kebab_case(router_name)
    human_name = ""
    fields = []
    model_name = ""
    if module_name != "":
        human_name = capitalize(module_name)
        model_name = pascal_case(module_name.lower())
        print("Create fields for "+model_name)
        finish = False
        while not finish:
            field = add_fields(len(fields), snake_case(default_router_name))
            if field["name"] != "":
                fields.append(field)
            else:
                finish = True
    model_struct = ""
    model_properties = []
    model_values = []
    model_json_values = []
    if len(fields) < 1:
        print("No model properties provided")
        exit(0)

    for i, f in enumerate(fields):
        model_properties.append(f["name"]+" "+f["type"]+" "+f["tag"]+f["comment"])
        model_values.append(f["name"]+": nil,")
        json_item = "\""+f["lower_name"]+"\": null"
        if i + 1 < len(fields):
            json_item = json_item + ","
        model_json_values.append(json_item)


    base_name = snake_case(module_name)
    package_name = lower_alphanumeric(base_name)

    model_struct = template_model.replace("{{MODEL_NAME}}", model_name)
    model_struct = model_struct.replace("{{MODEL_PROPERTIES}}", "\n".join(model_properties))
    model_struct = model_struct.replace("{{MODEL_NAME_PCASE}}", base_name)
    model_struct = model_struct.replace("{{MODEL_NAME_HUMAN}}", human_name)

    controller_dir = "app/controller/" + package_name
    model_path = "app/model/" + base_name + ".go"

    if not os.path.isfile(model_path):
        print("Creating models and endpoints ...")
        with open(model_path, "w") as fout:
            fout.write(model_struct)

        clean_code(model_path)

        create_dir(controller_dir)
        get_file = controller_dir + "/" + base_name + "_get.go"
        get_file_test = controller_dir + "/" + base_name + "_get_test.go"
        get_id_file = controller_dir + "/" + base_name + "_id_get.go"
        get_id_file_test = controller_dir + "/" + base_name + "_id_get_test.go"
        delete_file = controller_dir + "/" + base_name + "_delete.go"
        delete_file_test = controller_dir + "/" + base_name + "_delete_test.go"
        post_file = controller_dir + "/" + base_name + "_post.go"
        post_file_test = controller_dir + "/" + base_name + "_post_test.go"
        put_file = controller_dir + "/" + base_name + "_put.go"
        put_file_test = controller_dir + "/" + base_name + "_put_test.go"

        replacements = {
            "{{MODEL_NAME}}": model_name,
            "{{MODEL_NAME_SNAKE}}": base_name,
            "{{PACKAGE_NAME}}": package_name,
            "{{ROUTER_NAME}}": router_name,
            "{{MODEL_NAME_HUMAN}}": human_name,
            "{{MODEL_VALUES}}": "\n".join(model_values),
            "{{JSON_VALUES}}": "\n\t\t".join(model_json_values),
        }

        router_template = """
	// {{MODEL_NAME_HUMAN}}
	{{PACKAGE_NAME}}API := api.Group("/{{ROUTER_NAME}}")
	{{PACKAGE_NAME}}API.Use(middleware.TokenValidator())
	{{PACKAGE_NAME}}API.Post("/", {{PACKAGE_NAME}}.Post{{MODEL_NAME}})
	{{PACKAGE_NAME}}API.Get("/", {{PACKAGE_NAME}}.Get{{MODEL_NAME}})
	{{PACKAGE_NAME}}API.Put("/:id",  {{PACKAGE_NAME}}.Put{{MODEL_NAME}})
	{{PACKAGE_NAME}}API.Get("/:id", {{PACKAGE_NAME}}.Get{{MODEL_NAME}}ID)
	{{PACKAGE_NAME}}API.Delete("/:id", {{PACKAGE_NAME}}.Delete{{MODEL_NAME}})
        """

        migration_template = """
    // ModelMigrations models to automigrate
    var ModelMigrations = []interface{}{
        &model.{{MODEL_NAME}}{},
    }
        """

        with open(get_file, "w") as gin:
            content = template_get.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(get_id_file, "w") as gin:
            content = template_get_id.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(delete_file, "w") as gin:
            content = template_delete.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(post_file, "w") as gin:
            content = template_post.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(put_file, "w") as gin:
            content = template_put.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(get_file_test, "w") as gin:
            content = template_get_test.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(get_id_file_test, "w") as gin:
            content = template_get_id_test.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(post_file_test, "w") as gin:
            content = template_post_test.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(put_file_test, "w") as gin:
            content = template_put_test.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        with open(delete_file_test, "w") as gin:
            content = template_delete_test.split("---")[0]
            for k, v in replacements.items():
                content = content.replace(k, v)
            gin.write(content)

        # Router ...
        for k, v in replacements.items():
            router_template = router_template.replace(k, v)
            migration_template = migration_template.replace(k, v)

        generated = glob.glob("app/controller/"+package_name+"/*.go")
        for g in generated:
            clean_code(g)
        print("Completed!")
        print("")
        print("Add this code to router.go")
        print(router_template)
        print("")
        print("Add this code to migration.go")
        print(migration_template)
    else:
        print(model_path + " already exists")

create_endpoints()
