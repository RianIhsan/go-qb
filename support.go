package querybuilder

import (
	"fmt"
	"strings"
)

// DBType mendefinisikan tipe database: mysql/postgresql
type DBType int
type Param struct {
	Key   string
	Value interface{}
}

const (
	// DBPostgres merepresentasikan PostgreSQL
	DBPostgres DBType = 1

	// DBMySQL merepresentasikan MySQL
	DBMySQL DBType = 2
)

// Builder merepresentasikan query builder sederhana, terutama untuk mencegah injeksi SQL
type Builder struct {
	sb       *strings.Builder
	counter  int
	args     []interface{}
	bindType DBType
}

// New membuat query builder baru.
func New(dbType DBType, baseQuery string, data ...interface{}) *Builder {
	var sb strings.Builder

	b := &Builder{
		sb:       &sb,
		bindType: dbType,
	}
	b.addQuery(baseQuery, data...)
	return b
}

// AddQuery menambahkan query dengan format string dan data yang diberikan.
func (b *Builder) AddQuery(format string, data ...interface{}) {
	b.sb.WriteString(" ")
	b.addQuery(format, data...)
}

func (b *Builder) addQuery(format string, data ...interface{}) {
	b.args = append(b.args, data...)

	if b.bindType == DBPostgres {
		for i := 0; i < len(data); i++ {
			b.counter++
			format = strings.Replace(format, "?", fmt.Sprintf("$%d", b.counter), 1)
		}
	} else if b.bindType == DBMySQL {
		for i := 0; i < len(data); i++ {
			b.counter++
			format = strings.Replace(format, "?", fmt.Sprintf("?", b.counter), 1) // MySQL tetap menggunakan ?
		}
	}

	b.sb.WriteString(format)
}

// AddString menambahkan string mentah ke dalam query.
func (b *Builder) AddString(str string) {
	b.sb.WriteString(" " + str)
}

// Query mengembalikan query yang aman dengan placeholder yang diganti sesuai kebutuhan
func (b *Builder) Query() string {
	return b.sb.String()
}

// Args mengembalikan data argumen untuk query
func (b *Builder) Args() []interface{} {
	return b.args
}

// GenerateQuery membuat query dengan kondisi dinamis berdasarkan parameter.
func GenerateQuery(query string, fn func() []Param, dbType DBType) *Builder {
	args := fn()

	where := "WHERE "
	queryBuilder := New(dbType, query)
	for i := 0; i < len(args); i++ {
		if i != 0 {
			queryBuilder.AddQuery(" AND "+args[i].Key+" = ?", args[i].Value)
			continue
		}

		queryBuilder.AddQuery(where+args[i].Key+" = ?", args[i].Value)
	}

	return queryBuilder
}

// GenDynamicPlaceholderSQL menghasilkan string placeholder yang cocok untuk klausa IN.
func GenDynamicPlaceholderSQL(length int, dbType DBType) (str string) {
	placeholders := make([]string, length)
	for i := 0; i < length; i++ {
		if dbType == DBPostgres {
			placeholders[i] = fmt.Sprintf("$%d", i+1) // PostgreSQL
		} else {
			placeholders[i] = "?" // MySQL
		}
	}
	return strings.Join(placeholders, ", ")
}

// Metode Update untuk menangani operasi SQL pembaruan secara dinamis.
func (b *Builder) Update(table string, data ...Param) *Builder {
	b.sb.WriteString(fmt.Sprintf("UPDATE %s SET ", table))
	for i, param := range data {
		if i != 0 {
			b.sb.WriteString(", ")
		}
		b.AddQuery(fmt.Sprintf("%s = ?", param.Key), param.Value)
	}
	return b
}

// Metode Delete untuk menangani operasi penghapusan.
func (b *Builder) Delete(table string) *Builder {
	b.sb.WriteString(fmt.Sprintf("DELETE FROM %s ", table))
	return b
}

// Metode Select untuk menangani operasi pemilihan.
func (b *Builder) Select(columns string, table string) *Builder {
	b.sb.WriteString(fmt.Sprintf("SELECT %s FROM %s ", columns, table))
	return b
}
